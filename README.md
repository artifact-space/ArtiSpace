![](docs/assets/20250205_124228_logo_8.png)

# ArtiSpace

ArtiSpace is an open-source artifact registry that supports multiple package types(under development).<hr/>

## Why ArtiSpace?

The ArtiSpace is an open source focused artifact registry that intends to support popular package types.

ArtiSpace is an open-source artifact registry designed to support popular package types.

Unlike existing artifact registries, ArtiSpace prioritizes **simplicity, ease of deployment, and modern technology** . Our goal is to provide a lightweight yet powerful solution that is **developer-friendly, scalable, and efficient**

Compared to already available artifact registries, the ArtiSpace focues more on easy to deploy, simplicity and is being built on top of latest technologies.

## Architecture

ArtiSpace features a simple and layered architecture designed for clarity and scalability.
At a high level, it consists of four layers:

![](docs/assets/20250205_124318_ArtiSpace_Architecture.drawio.png)

1. API Layer – Handles incoming requests and exposes RESTful endpoints.
2. Application Layer – Contains ArtiSpace's core logic.
3. Metadata Layer – Manages package metadata and indexing.
4. Storage Layer – Stores artifacts efficiently.

ArtiSpace is designed with user-friendly interfaces to simplify management for both developers and administrators.

## Technologies

* **Backend (Core Logic):** Golang
* **Frontend (UI):** React
* **Databases:** RDBMS for metadata storage. Initially, we are developing a**proof of concept (PoC)** with**SQLite** , with plans to extend support to**PostgreSQL** and other RDBMS solutions.
* **Storage:** The PoC will use a**local file system** , with future support planned for**S3, Azure Blob Storage, Google Cloud Storage, and NFS**
* App Logic/Core/Backend - Golang
* UI - React
* DBs : RDBMS (used for meta data storage) and intially, We intened to develop PoC with SQLite and extends to other RDMBS such as PostgreSQL later.
* Storage: PoC will be developed with Local File system. Later, We planned to support other storage backends such as S3 / Azure Blob Storage / Google Cloud Storage and NFS

---

## Docker Support in ArtiSpace

ArtiSpace will support the **Docker V2 API**. The `docker push` command involves multiple HTTP requests to the registry. Understanding these API calls is crucial for designing the database schema and algorithms to handle image pushes flawlessly.

### Docker Push Workflow

**Example Command**:
`docker push localhost:5000/john/busybox:0.1`

The above is an example docker push command. In the above command,

* localhost:5000 - DockerHub / Registery's address
* john - Docker Name Space / If no namespace is specified, docker's default = library
* busybox - Docker Repository
* 0.1 - Docker Tag

**Version Check**

As soon as the above command is executed, the Docker daemon will send a request to check if the Docker V2 API is supported.

`GET https://localhost:5000/v2`

If the above request receives an `HTTP 200` status code as a response, the Docker daemon will assume that the registry supports the Docker V2 API and continue pushing the image.

**Pushing Docker Container Image Config and Layers as blobs**

After the version check, the Docker daemon will upload Docker Container Image Config and Layers as blobs.

* Container Image Config -> Blob
* Layer 1 -> Blob
* Layer 2 -> Blob
* Layer 3 -> Blob

There are two ways to upload blobs:

* Monolithic Blob Uploads: The entire layer will be uploaded as one chunk.
* Chunked Blob Uploads: Layers will be split into multiple chunks and uploaded separately in HTTP requests.

**Checking for Existing Layer**

The Docker daemon will send HEAD requests to check the existence of layers and the container image config.

`HEAD https://localhost:5000/v2/john/busybox/sha256:0ef2e08ed3fabfc44002ccb846c4f2416a2135affc3ce39538834059606f32dd`

If the response receives a 200 status code, the Docker daemon won't upload the layer. Otherwise, it will start uploading the layer/container image config as a blob.

**Uploading a blob**

Before uploading a blob, the Docker daemon must obtain a URL from the registry. The obtained URL will be used (in multiple requests if chunked blob uploads) until the Docker daemon finishes uploading the entire layer/blob.

The response to the above request should have a Location header with a valid URL where the blob can be uploaded. Also, the status code must be 202.

An example URL:
`https://localhost:5000/v2/john/busybox/blob/uploads/1a39f753-36b0-4cec-a098-a55118c2e435`

In ArtiSpace, we refer to 1a39f753-36b0-4cec-a098-a55118c2e435 as the blob upload session ID.

The Docker daemon can decide how to upload the blobs: Monolithic Upload or Chunked Upload.

**Monolithic Blob Upload**

For monolithic blob uploads, the following PUT request will be sent to upload the entire layer or container image config:

` PUT https://localhost:5000/v2/john/busybox/blob/uploads/1a39f753-36b0-4cec-a098-a55118c2e435?digest=sha256:0ef2e08ed3fabfc44002ccb846c4f2416a2135affc3ce39538834059606f32dd`

**Chunked Blob Upload**

In chunked uploads, PATCH requests will be used to upload chunks except for the final one. The final chunk will be uploaded with a PUT request.

Example:

```
Chunk 1: PATCH https://localhost:5000/v2/john/busybox/blob/uploads/1a39f753-36b0-4cec-a098-a55118c2e435 
Chunk 2: PATCH https://localhost:5000/v2/john/busybox/blob/uploads/1a39f753-36b0-4cec-a098-a55118c2e435
....

Final Chunk: PUT https://localhost:5000/v2/john/busybox/blob/uploads/1a39f753-36b0-4cec-a098-a55118c2e435?digest=sha256:0ef2e08ed3fabfc44002ccb846c4f2416a2135affc3ce39538834059606f32dd

```

The final chunk always includes the query parameter digest with the digest of the layer/container image config.

After uploading the container image config and layers, the image manifest will be uploaded.

**Manifest Upload**

` PUT https://localhost:5000/v2/john/busybox/manifests/v1`