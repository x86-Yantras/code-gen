## Open api code generator
---
An opinionated code geneartor for API driven development
### Featues Supported
---
- [x] Languages Supported
    - [x] Node.js
    - [ ] Golang
- [x] Project Initialization
- [x] Service Generation
    - [x] Required Validation generation
    - [x] Custom error generation
    - [x] Test generation with storage mock
- [x] Adapter Generation
    - [x] Http Generation
    - [x] Storage Generation

### Usage
---
> code-gen [path-to-api-spec] [language] [command] [serviceName(optional)]

### Commands list
--- 
- `init`: Initializes project/adds new updates
- `services`: generate services
- `http`: generates http layer
- `storage`: generates storage layer

### Notes
---
- Please run npm i after running init command
- Passing optional serviceName generates only for that service
- Sevice name should match specs.tag[0] in api spec file

### How to Start
---
#### Generate Project
> code-gen [path-to-api-spec] [language] init
- This step generates a bare project with server, loggers, config and request helpers.
- The server file can be changed, if custom server implementation is needed

#### Generate Services
> code-gen [path-to-api-spec] [language] services
- Things to know:
    - The service functions are grouped according to first tag in tags in spec file
    - Service/handler Operations(like creteUser) will be generated using OperationId, so it is must in api spec
    - `types` file will be over-written when running commands in future.
    - Custom Validations should be kept in the service operation file and not in types
    - Custom errors will be generated if there is `errors` in api spec components
   
    - Custom error codes will be generated if there is `errCodes` in api spec components
    ```
    components:
      errors:
        - name: ValidationError
          httpStatusCode: 400
      errCodes:
        - name: ValueRequired
          code: "_value_required"
    ```

    - For automatic required validations use `required` feature of openapi spec
    ```
    components:
      schemas:
        CreateUserPayload:
          type: object
          required:
          - name
          - gender
          - birthday
          properties:
            name: 
              type: string
            gender:
              type: string
            birthday:
              type: string
    ```
#### Generate Http Adapters
> code-gen [path-to-api-spec] [language] http
- Things to know:
    - The handlers are grouped according to first tag in tags in spec file
    - `routes` file will be over-written when running commands in future.

### Generate Storage Adapters
> code-gen [path-to-api-spec] [language] storage
- This step generates storage, default mongo, adapter
- Things to know:
    - The storage are grouped according to first tag in tags in spec file
    - At least one component/schema with entity name is must for generating the storage
    ```
    components:
      schemas:
        User:
        type: object
        properties:
            name: 
            type: string
    ```

### Generated Project structure
    tree
    |-- src // contains all the source code
    |   |-- services // conatins all the services
    |   |   |-- [servicename] // service directory for the domain
    |   |   |   |   |-- adapters // adapters for the service
    |   |   |   |   |   |-- api // all the api adapters
    |   |   |   |   |   |   |-- http
    |   |   |   |   |   |   |-- async
    |   |   |   |   |   |-- storage // all the storage adapters
    |   |   |   |   |   |   |-- mongo
    |   |   |   |   |   |   |-- postgresql // not implemented
    |   |-- lib // contains all the utilities and libraries
    |   |   |-- context // creates context for request cycle
    |   |   |-- errors // custom error implementations and codes
    |   |   |-- middlewares // request middlewares
    |   |   |   |-- logRequest
    |   |   |   |-- requestId
    |   |   |-- test // test helpers
    |   |   |-- db // db helpers and connectors
    |   |   |-- logger
    |   |-- config
    |-- [indexFile] // contains server and init code
    |-- [packageManagerFile]
    |-- Readme.md
    |-- .gitignore
    