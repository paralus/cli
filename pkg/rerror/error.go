/*
Define errors here. For example, errors for when we cannot retrieve
a resource, or a timeout error.
*/
package rerror

type ResourceNotFound struct {
	Type string
	Name string
}

type GenericError struct {
	Message string
}

func (e ResourceNotFound) Error() string {
	return "resource " + e.Name + " of type " + e.Type + " not found"
}

func (e GenericError) Error() string {
	return e.Message
}

type CrudErr struct {
	Type string
	Name string
	Op   string
}

func (e CrudErr) Error() string {
	return "could not complete operation " + e.Op + " on resource " + e.Name + " of type " + e.Type
}

type FileUploadErr struct {
	FileName string
}

func (e FileUploadErr) Error() string {
	return "could not upload file " + e.FileName
}
