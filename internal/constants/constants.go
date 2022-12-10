package constants

// Dirs
const TemplatesDir = "templates"
const ConfigDir = "./config"
const ServiceDirPlaceholder = `servicename`

// Files
const TemplateExtension = ".tmpl"
const ServiceFilePlaceholder = `service`
const HandlerPlaceHolder = `handler`
const TypesFile = "types"
const TypesValidatorFile = "types_validation"

// Messages
const ProjectBuiltMsg = `Built %s, don't forget to commit`
const UndefinedCommandMsg = `Cannot parse command %s`
const TagMissingErr = `tags missing in method in spec file, please check example file`
const TagDifferentErr = `different tags: %s under same path: %s`

// HTTP
const PayloadLimit = "limit"
const PayloadOffset = "offset"
const ContentJson = "application/json"
const APIHTTPAdapter = "api/http"

// Storage
const Storage = "storage"
const StorageEntityPlaceholder = "entity"
const StorageEntityOpPlaceholder = "entity_op"

// To-do(rolorin): move this to cli ops
const DefaultStorage = "mongo"

// misc
const Service = "service"
const Errors = "errors"
const ErrCodes = "errCodes"
