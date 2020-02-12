# Wordpress Elastic Search

![go version](https://img.shields.io/github/go-mod/go-version/igorgottschalg/wordpress-elasticsearch-index)
![docker build](https://img.shields.io/docker/build/gottschalg/wordpress-elasticsearch-index)

This microservice has a proposal to receive a wordpress data and index on ElasticSearch Database.

### How to use
In wordpress save post hook, implement or use a plugin to send data to this microservice, with the bellow structure:
```
{
    "id": Interger,
    "name": String,
    "content": String,
    "image": String,
    "url": String,
    "post_type": String,
    "keywords": Array of String
}
```
