App Dependency Discoverer
=========================

This application can discover dependencies of application stack with provided root GUID. Application in CF can be bound to services and user provided services. When you need to link to other application, use user provided service with url field filled with other application address (starting with scheme), e.g.
```
cf cups <hostname>-ups -p url
url> http://<hostname>.<domain>
```
Applications linked must be in the same space.


This is an app which:
 
 * inquiry CF for dependencies of application provided in rootGUID, 
 * constructs dependency tree,
 * check for cycles,
 * if they exist, application stack cannot be cloned (there is no order it could be spawned),
 * if dependency tree is direct acyclic graph, application returns a list of components from specified application stack in reversed topological order (in which they can be cloned): guid, name, type, list of dependent components and other pieces of information useful when cloning
   

### Idea behind

This app can be used by new [application broker](https://github.com/trustedanalytics/application-broker/) to retrieve list of components which should be cloned when spawning application stack based on existing stack.

GET /v1/discover/< rootGUID >
```
Example response body:
[
  {
    "GUID": "b12e08f1-0329-471e-9cc7-9a26bb24b072",
    "name": "mycdh",
    "type": "Service",
    "dependencyOf": [
      "ce44ee69-32b4-4e6f-a952-a10b0d522021",
      "90492a34-1f00-43b5-bcec-828456d8981a"
    ],
    "clone": true
  },
  {
    "GUID": "ca47e7b4-cf3f-443b-ac40-a5c99e36b232",
    "name": "myhdfs",
    "type": "Service",
    "dependencyOf": [
      "ce44ee69-32b4-4e6f-a952-a10b0d522021",
      "90492a34-1f00-43b5-bcec-828456d8981a"
    ],
    "clone": true
  },
  {
    "GUID": "f09c5d58-34fb-4a69-9658-062314be9711",
    "name": "myrabbit",
    "type": "Service",
    "dependencyOf": [
      "ce44ee69-32b4-4e6f-a952-a10b0d522021"
    ],
    "clone": true
  },
  {
    "GUID": "90492a34-1f00-43b5-bcec-828456d8981a",
    "name": "app1",
    "type": "Application",
    "dependencyOf": [
      "ff8e11cf-de1a-4c27-abdf-a827d447e74a"
    ],
    "clone": true
  },
  {
    "GUID": "ff8e11cf-de1a-4c27-abdf-a827d447e74a",
    "name": "app1-ups",
    "type": "User provided service",
    "dependencyOf": [
      "ce44ee69-32b4-4e6f-a952-a10b0d522021"
    ],
    "clone": true
  },
  {
    "GUID": "ce44ee69-32b4-4e6f-a952-a10b0d522021",
    "name": "toplvlapp1",
    "type": "Application",
    "dependencyOf": [],
    "clone": true
  }
]
```

### IDE
We recommend using [IntelliJ IDEA](https://www.jetbrains.com/idea/) as IDE with [golang plugin](https://github.com/go-lang-plugin-org/go-lang-idea-plugin). To apply formatting automatically on every save you may use go-fmt with [File Watcher plugin](http://www.idmworks.com/blog/entry/automatically-calling-go-fmt-from-intellij).


Tips
-----------------------

### Golang tips

Developing golang apps requires you store all dependencies (Godeps) in separate directory. They shall be placed in source control.

```
godep save ./...
```

Command above places all dependencies from `$GOPATH`, your app uses, in Godeps and writes its versions to Godeps/Godeps.json file.

### TODO

* Allow to reuse parts of stack in next clones