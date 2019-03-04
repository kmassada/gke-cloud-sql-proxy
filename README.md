# CloudSQl on k8s

## Create new instance sql

```shell
export CLOUD_SQL_NAME=instance1
gcloud sql instances create $CLOUD_SQL_NAME \
    --tier=db-n1-standard-2 \
    --region=us-west2
gcloud services enable sqladmin.googleapis.com
```

## Download instance name

```shell
export INSTANCE_CONNECTION_NAME=`gcloud sql instances describe $CLOUD_SQL_NAME --format='value(connectionName)'`
kubectl create secret generic cloudsql-instance-connection \
 --from-literal connection=$INSTANCE_CONNECTION_NAME \
 --from-literal name=$CLOUD_SQL_NAME
```

## Create SA

```shell
gcloud iam service-accounts create $CLOUD_SQL_NAME --display-name "Cloud SQL Service Account"
export CLOUD_SQL_EMAIL=`gcloud iam service-accounts list --format='value(email)' --filter='displayName:Cloud SQL Service Account'`

export PROJECT_ID=`gcloud config get-value project`

gcloud projects add-iam-policy-binding $PROJECT_ID --member=serviceAccount:${CLOUD_SQL_EMAIL} --role=roles/cloudsql.editor

```

## Download Connection Credentials

```shell
gcloud iam service-accounts keys create \
    /home/$USER/$CLOUD_SQL_NAME-key.json \
    --iam-account $CLOUD_SQL_EMAIL

kubectl create secret generic cloudsql-instance-credentials --from-file /home/$USER/$CLOUD_SQL_NAME-key.json
```

```shell
date +%s | sha256sum | base64 | head -c 12 > file
gcloud sql users create admin \
   --host=% --instance=$CLOUD_SQL_NAME --password=`cat file`
kubectl create secret generic cloudsql-db-credentials \
 --from-literal username=admin \
 --from-literal password=`cat file`
 ```

## Verify

creds

```console
$ kubectl get secrets
NAME                            TYPE                                  DATA      AGE
cloudsql-db-credentials         Opaque                                2         4m
cloudsql-instance-connection    Opaque                                1         11m
cloudsql-instance-credentials   Opaque                                1         7m
```

cloud sql proxy

```shell
wget https://dl.google.com/cloudsql/cloud_sql_proxy.linux.amd64 -O cloud_sql_proxy
chmod +x cloud_sql_proxy

./cloud_sql_proxy \
    -instances=$INSTANCE_CONNECTION_NAME=tcp:3306 \
    -credential_file=/home/$USER/$CLOUD_SQL_NAME-key.json &

> mysql -h 127.0.0.1 -u root --protocol=tcp
```

## Wordpress

### Sub and Apply

```shell
envsubst < mysql_wordpress_deployment.yaml > deployment.yaml
kubectl apply -f deployment.yaml
```

## cloudsql-client

### build

```shell
export _REPO_PREFIX=makz-labs
export APPLICATION=gke-cloud-sql-proxy
export REPO_NAME=$APPLICATION
gcloud builds submit \
    --substitutions _REPO_PREFIX=$_REPO_PREFIX,REPO_NAME=$REPO_NAME,TAG_NAME=cli,SHORT_SHA=clisha \
    --config cloudbuild.yaml
```

### Apply

```shell
envsubst < mysql_grpc_deployment.yaml > deployment.yaml
kubectl apply -f deployment.yaml
```

tables need to be created for this example

```sql
CREATE DATABASE cloudsqlclient;
USE cloudsqlclient;
DROP TABLE squareNum;
CREATE TABLE squareNum (
    number int,
    squareNumber int,
    PRIMARY KEY( number )
);
```