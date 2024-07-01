# Flock: A Framework for Deploying On-Demand Distributed-Trust

Flock is a framework for deploying on-demand distributed-trust systems. This implementation contains our system as described in our [paper](https://eprint.iacr.org/2024/892.pdf) accepted into OSDI 2024.

**WARNING:** This is an academic proof-of-concept prototype and has not received careful code review. This implementation is NOT ready for production use.

## Table of Contents
- [Artifact Evaluation](#artifact-evaluation)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
  - [Generate User Certificates](#generate-user-certificates)
  - [Building Flock's Docker Image](#building-flocks-docker-image)
- [Running Flock](#running-flock)
  - [Baseline Setup](#baseline-setup)
  - [Example Applications](#example-applications)
- [Serverless Deployment](#serverless-deployment)
  - [Amazon Lambda](#deploying-amazon-lambda)
  - [Google Cloud Run](#deploying-to-google-cloud-run)
- [Building the Flock Relay](#build-the-flock-relay)
- [Citation](#citation)

## Artifact Evaluation
Flock has received all three badges (Available, Functional, Reproduced) from OSDI 2024's artifact evaluation committee. The documentation for our artifact evaluation is archived in the [/artifacts](/artifacts/README.md) folder. The previous repo for Flock's artifact evaluation is archived and can be found [here](https://github.com/flock-org/flock-artifact)

## Prerequisites
- Docker
- Python 3
- Bazel (for PIR application)
- Cloud provider accounts (AWS, GCP, Azure) for distributed setup

## Getting Started

### Generate User Certificates
Certificates are generated using the tools in the `/relay` folder. Follow these steps:

```bash
cd /relay
make build
./bin/fr-adm create party --name 0 --user user1
./bin/fr-adm create party --name 1 --user user1
./bin/fr-adm create party --name 2 --user user1
```

Certificates will be generated in the `certs` folder.

### Building Flock's Docker Image
You can either pull the pre-built image or build it yourself:

```bash
# Pull pre-built image
docker pull sijuntan/flock:ubuntu

# Or build it yourself
docker build --platform linux/amd64 -t flock -f Dockerfile.ubuntu .
```

## Running Flock

### Baseline Setup
To run Flock in a baseline setup with 3 VMs (ideally one on each cloud provider), execute the following command on each VM:

```bash
# On VM 1 (e.g. AWS)
sudo docker run -p 443:443 -p 5000-5100:5000-5100 \
  -e PARTY_CERT="$(cat certs/user1/0/cert.pem)" \
  -e PARTY_KEY="$(cat certs/user1/0/key.pem)" \
  sijuntan/flock:ubuntu python3 handler.py -s local

# On VM 2 (e.g. GCP)
sudo docker run -p 443:443 -p 5000-5100:5000-5100 \
  -e PARTY_CERT="$(cat certs/user1/1/cert.pem)" \
  -e PARTY_KEY="$(cat certs/user1/1/key.pem)" \
  sijuntan/flock:ubuntu python3 handler.py -s local

# On VM 3 (e.g. Azure)
sudo docker run -p 443:443 -p 5000-5100:5000-5100 \
  -e PARTY_CERT="$(cat certs/user1/2/cert.pem)" \
  -e PARTY_KEY="$(cat certs/user1/2/key.pem)" \
  sijuntan/flock:ubuntu python3 handler.py -s local
```

### Example Applications
All client-side invocation codes are in the `/client` folder. Before running examples, update the IP addresses in `/client/config.py`.

```bash
cd /client

# Secret Recovery
python3 invoke.py baseline sharding_setup 10
python3 invoke.py baseline sharding 10

# File Decryption
python3 invoke.py baseline aes_setup 1
python3 invoke.py baseline aes_encrypt 1

# Signing
python3 one_time_setup
python3 invoke.py baseline signing_keygen 10
python3 invoke.py baseline signing_sign 10

# Data Freshness
python3 one_time_setup
python3 invoke.py baseline freshness_store_file 10
python3 invoke.py baseline freshness_retrieve_file 10

# PIR (Private Information Retrieval)
# First, build the necessary executables
cd /applications/pir
bazel build //:client_gen_pir_requests_bin
bazel build //:client_handle_pir_responses_bin
cd /client
python3 invoke.py baseline pir_setup 10
python3 invoke.py baseline pir 10
```

## Serverless Deployment

### Deploying Amazon Lambda
AWS Lambda requires the Docker image to be in AWS ECR. Follow these steps to pull Flock's image for AWS Lambda and deploy it there:

```bash
docker pull sijuntan/flock:lambda
docker tag sijuntan/flock:lambda <AWS-ECR-REGISTRY-ADDRESS>/flock:lambda
docker push <AWS-ECR-REGISTRY-ADDRESS>/flock:lambda
python3 ./deploy/deploy_aws.py --image_uri <AWS-ECR-REGISTRY-ADDRESS>/flock:lambda
```

### Deploying to Google Cloud Run
Deploy to Google Cloud Run with:

```bash
gcloud run deploy --image sijuntan/flock:ubuntu
```

## Build the Flock Relay
The Flock relay is located in [./relay/](./relay/). Build it with:

```bash
cd relay/
make build
```

The relay binary will be located in `./relay/bin/`. More details about its deployment can be found [here](/relay/README.md).

## Citation
If you use Flock in your research, please cite our paper:

```bibtex
@inproceedings{kaviani2024flock,
  title={Flock: A Framework for Deploying On-Demand Distributed Trust},
  author={Kaviani, Darya and Tan, Sijun and Kannan, Pravein Govindan and Poda, Raluca Ada},
  booktitle={USENIX Symposium on Operating Systems Design and Implementation},
  year={2024}
}
```
