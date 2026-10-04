package storage

type Storage interface{
CreateStudent(name string, email string, age int64) (int64, error)
}

