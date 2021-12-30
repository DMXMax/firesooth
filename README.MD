## CLI to work with firestore

## Setup

### Connecction to the Emulator

* Make sure emulator is running
* ```export FIRESTORE_EMULATOR_HOST=<ip>:<port>```

### Config File

* Located in home directory. Currently named ".config"
* Right now, it only holds the project name, which doesn't matter with an emulator, but will matter when talking to firestore

### Connecting to cloud firestore

if FIRESTORE_EMULATOR_HOST is set, the clients connect to the emulator host, regardless of project name.  

So far, so good.

We can add documents.
```Tree``` shows the collection tree. Maybe we make it so we can add a root path.

