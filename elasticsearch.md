## Elastic Search
To export to an Elastic Search database, populate `config.yaml`:
```
destination: elastic

elasticsearch:
  uri: http://127.0.0.1:9200
```

### Zinc using Amazon S3 buckets
A very quick and easy way to get up and running an Elastic Search type database, using Amazon AWS S3 buckets for storage is to use [Zinc](https://github.com/zinclabs/zinc), which is a single (Go) binary implementation of Elastic Search. It provides a nice web user interface, similar to Kibana, but self contained in the same binary as the Elastic Search engine.

In `config.yaml` a username and password is required. Set `use-zinc-aws-s3` to true. 
```
elasticsearch:
  uri: http://127.0.0.1:4080
  username: admin
  password: admin

  use-zinc-aws-s3: true
  no-duplicates: true
```
Note: `no-duplicates` is used to reduce the volume of data stored. Each document index is considered unique by the ID being a hash of the domain name and repository name (user/repo). Older documents with the same ID will be updated with newer commits as the arrive.

Example:
```
wget https://github.com/zinclabs/zinc/releases/download/v0.3.2/zinc_0.3.2_Linux_x86_64.tar.gz
tar xf zinc_0.3.2_Linux_x86_64.tar.gz

export ZINC_FIRST_ADMIN_PASSWORD=admin
export ZINC_FIRST_ADMIN_USER=admin
export AWS_ACCESS_KEY_ID=XXXXXXXXXXXXXX
export AWS_SECRET_ACCESS_KEY=XXXXXXXXXXXXXX
export AWS_REGION=us-east-1
export ZINC_S3_BUCKET=cstream

./zinc &
./commit-stream &
```

Then browsing to http://127.0.0.1/4080:
![Image](https://i.imgur.com/kNBKSkP.png)
