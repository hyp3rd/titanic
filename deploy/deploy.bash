#!/bin/bash
set -Eeuo pipefail

traperr() {
  echo "ERROR: ${BASH_SOURCE[1]} at about line ${BASH_LINENO[0]}"
}

set -o errtrace
trap traperr ERR

validate_env () {
    if [[ -z ${PROJECT_ID+x} ]] || [[ -z ${REGION+x} ]]; then
      echo "To run this deployment you need to export PROJECT_ID and REGION as follows:
      export REGION=<region e.g. europe-west1>
      export PROJECT_ID=<project name e.g. hyperd-titanic-api>";
      exit 1
    fi
}

gcloud_setup () {
    # implicitly enable the apis we'll use
  gcloud services enable storage-api.googleapis.com
  gcloud services enable cloudresourcemanager.googleapis.com
  gcloud services enable compute.googleapis.com
  gcloud services enable container.googleapis.com
  gcloud services enable iam.googleapis.com

  # create a service account in the project, give it permissions, and obtain a key for terraform to use
  if gcloud iam service-accounts list | grep "terraform@$PROJECT_ID.iam.gserviceaccount.com" | awk '{print $1}'; then
    echo "service-accounts: terraform@$PROJECT_ID.iam.gserviceaccount.com"
  else
    gcloud iam service-accounts create terraform --display-name "terraform"
    gcloud projects add-iam-policy-binding $PROJECT_ID --member "serviceAccount:terraform@$PROJECT_ID.iam.gserviceaccount.com" --role "roles/owner"
    gcloud iam service-accounts keys create key.json --iam-account terraform@$PROJECT_ID.iam.gserviceaccount.com
  fi

  # make the key available in this session
  export GOOGLE_APPLICATION_CREDENTIALS="$PWD/key.json"

  # implicitly install other gcloud components we'll need to run the rest
  gcloud components install -q gsutil kubectl docker-credential-gcr
}

terraform_deployment () {
  # create the shared state bucket for terraform to save it being persisted locally / allow other people to run the tooling
  gsutil mb -l $REGION gs://$PROJECT_ID-terraform-state || :

  # initialise terraform state and providers
  ./utils/terraform init -backend-config=bucket=$PROJECT_ID-terraform-state

  # apply the .tf manifests i.e. create the infrastructure and the cluster
  ./utils/terraform apply
}

build_docker_images () {
  # get some creds for the cluster to use for the next bit
  gcloud container clusters get-credentials titanic-api-cluster --region $REGION

  # explicitly specify the cluster in case we have others configured
  kubectl config use-context gke_${PROJECT_ID}_${REGION}_titanic-api-cluster

  # build the titanic-api image with our modifications (see Dockerfile) and tag for private GCR
  docker build --file ../docker/Dockerfile -t gcr.io/$PROJECT_ID/titanic-api .

  # configure pushing to private GCR, and push our image
  gcloud auth configure-docker -q
  docker push gcr.io/$PROJECT_ID/titanic-api
}

init_cluster_resources () {
  # apply the kubernetes manifests to provision the namespaces we need
  kubectl create --save-config -f k8s/namespaces/ || :

  # apply the kubernetes manifests which are declarative:
  # generic k8s
  kubectl apply -f k8s/
}

deploy_secrets () {
  # create a private key and a self signed certificate (remember that old skool 2048 bit as Google load balancers don't like the stronger RSA-4096)
  if [[ ! -f $(pwd)/tls/titanic-api.hyperd.sh.key ]] || [[ ! -f $(pwd)/tls/titanic-api.hyperd.sh.crt ]]; then
    openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 \
    -subj "/C=NL/ST=Amsterdam/L=Amsterdam/O=hyperd/CN=titanic-api.hyperd.sh" \
    -keyout tls/titanic-api.hyperd.sh.key -out tls/titanic-api.hyperd.sh.crt

    # make sure to remove the secret titanic-api-tls, if exists
    kubectl -n hcs delete secret titanic-api-tls || :
  fi

  # create a strong Diffie-Hellman group, used in negotiating Perfect Forward Secrecy with clients
  if [[ ! -f $(pwd)/tls/dhparam.pem ]]; then
    openssl dhparam -out tls/dhparam.pem 2048

    # make sure to remove the secret titanic-api-tls-dhparam, if exists
    kubectl -n hcs delete secret titanic-api-tls-dhparam || :
  fi

  # use an imperative command to create a kubernetes secret from this key that can be used with the GCE ingress
  tls_secret=$(kubectl -n hcs get secrets | awk '{print $1}' | awk -F, '$1 == V' V="titanic-api-tls")
  if [[ -z "$tls_secret" ]]; then
    kubectl -n hcs create secret tls titanic-api-tls --key tls/titanic-api.hyperd.sh.key --cert tls/titanic-api.hyperd.sh.crt
  fi

  tls_dhparam_secret=$(kubectl -n hcs get secrets | awk '{print $1}' | awk -F, '$1 == V' V="titanic-api-tls-dhparam")
  if [[ -z "$tls_dhparam_secret" ]]; then
    kubectl -n hcs create secret generic titanic-api-tls-dhparam --from-file=tls/dhparam.pem
  fi
}

init() {

  # check that PROJECT_ID and REGION are exported in the current shell
  validate_env

  # run the initial gcloud setup
  gcloud_setup

  # deploy the k8s cluster with terraform
  terraform_deployment

  # build the API docker images and push em to a GCR private registry
  build_docker_images

  # deploy the necessary initial resources to our GKE cluster
  init_cluster_resources

  # deploy the secrets to run the API app
  deploy_secrets
}

init

deploy_titanic_api () {
  # nginx ingress controller
  kubectl apply -f k8s/ingress-nginx/

  # deploy the api
  kubectl apply -f k8s/titanic-api/
}


# TODO: Automate the commented steps, properly waiting for the resources
# to be provisioned
deploy_cockroachdb () {
   # cockroachdb deployment init
  kubectl apply -f k8s/cockroachdb/

  # kubectl wait --for=condition=complete --timeout=60s pods --all || :

  # kubectl get csr

  # kubectl certificate approve default.node.cockroachdb-0

  # kubectl apply -f k8s/cockroachdb/cockroachdb-cluster-init/cluster-init-secure.yaml

  # kubectl --for=condition=complete --timeout=60s pods --all || :

  # kubectl get csr

  # kubectl certificate approve default.client.root

  # kubectl get job cluster-init-secure

  # kubectl get pods
}

deploy_titanic_api

deploy_cockroachdb