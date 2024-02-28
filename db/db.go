package db

import (
    "log"
    "database/sql"

    "git.fuwafuwa.moe/x3/ngfshare/id"
    "git.fuwafuwa.moe/x3/ngfshare/model"

    _ "github.com/mattn/go-sqlite3"
)

type DB struct {
    ctx *sql.DB
}

var Db *DB

func createTables(ctx *sql.DB) error {
    _, err := ctx.Exec(`
    CREATE TABLE IF NOT EXISTS keys (
        key TEXT PRIMARY KEY UNIQUE
    );

    CREATE TABLE IF NOT EXISTS files (
        id TEXT PRIMARY KEY UNIQUE,
        filename TEXT,
        size INTEGER,
        content_type TEXT,
        upload_time INTEGER,
        sha1sum TEXT UNIQUE,
        uploadKey TEXT,
        FOREIGN KEY (uploadKey)
            REFERENCES keys(key)
    );
    `)

    return err
}

func Open(path string) (*DB, error) {
    log.Printf("Opening DB at '%s'\n", path)

    ctx, err := sql.Open("sqlite3", path)
    if err != nil {
        log.Println("Failed to open DB", err)
        return nil, err
    }

    err = createTables(ctx)
    if err != nil {
        log.Println("Failed to create tables", err)
        return nil, err
    }

    Db = &DB{
        ctx: ctx,
    }
    return Db, nil
}

func (db *DB) Close() error {
    log.Println("Closing database...")
    err := db.ctx.Close()
    Db = nil
    return err
}

func (db *DB) IsAuthKeyExists(key string) bool {
    row := db.ctx.QueryRow(`
        SELECT EXISTS(SELECT 1 FROM keys WHERE key = ?);
    `, key)
    var ex int
    row.Scan(&ex)
    return ex == 1
}

func (db *DB) CreateNewAuthKey() (string, error) {
    var err error
    key := id.GenAuthKey()

    for i := 0; i < 10; i++ {
        _, err = db.ctx.Exec(`
            INSERT INTO keys (
                key
            )
            VALUES (?)
        `, key)

        if err == nil {
            break
        }
        /*
        errorCode := err.(sqll.Error).Code
        log.Println(errorCode)
        log.Println(err)
        if errorCode != sqll.ErrConstraint {
            log.Println("HERERERE")
            break
        }
        */

        // Try again
        key = id.GenAuthKey()
    }
    return key, err
}

func (db *DB) GetFileBySha1(sum string) (model.File, bool) {
    row := db.ctx.QueryRow(`
        SELECT
            id, filename, size, content_type, upload_time, uploadKey
        FROM
            files
        WHERE
            sha1sum = ?
    ;`, sum)

    f := model.File{Sha1Sum: sum}
    err := row.Scan(&f.Id, &f.Filename, &f.Size, &f.ContentType, &f.UploadTime, &f.UploadKey)
    //log.Println("db: ", err)
    return f, err == nil
}

func (db *DB) GetFileById(id string) (model.File, error) {
    row := db.ctx.QueryRow(`
        SELECT
            filename, size, content_type, upload_time, sha1sum, uploadKey
        FROM
            files
        WHERE
            id = ?
    ;`, id)

    f := model.File{Id: id}
    err := row.Scan(&f.Filename, &f.Size, &f.ContentType, &f.UploadTime, &f.Sha1Sum, &f.UploadKey)
    //log.Println("db: ", err)
    return f, err
}

func (db *DB) GetFilesByAuthKey(key string) ([]model.File, error) {
    lst := make([]model.File, 0, 16)
    rows, err := db.ctx.Query(`
        SELECT
            id, filename, size, content_type, upload_time, sha1sum
        FROM
            files
        WHERE
            uploadKey = ?
        ORDER BY
            upload_time DESC
    ;`, key)
    if err != nil {
        return lst, err
    }
    defer rows.Close()

    for rows.Next() {
        f := model.File{
            UploadKey: key,
        }
        err = rows.Scan(&f.Id, &f.Filename, &f.Size, &f.ContentType, &f.UploadTime, &f.Sha1Sum)
        if err != nil {
            return lst, err
        }
        lst = append(lst, f)
    }

    return lst, nil
}

func (db *DB) InsertFile(filename string, size int64, content_type, sha1sum, uploadKey string) (string, *sql.Tx, error) {
    var err error
    tx, err := db.ctx.Begin()
    if err != nil {
        log.Println("Cannot start Tx", err)
        return "", nil, err
    }
    fId := id.GenFileId()

    for i := 0; i < 10; i++ {

        _, err = tx.Exec(`
            INSERT INTO files (
                id,
                size,
                filename,
                content_type,
                upload_time,
                sha1sum,
                uploadKey
            )
            VALUES(?,?,?,?,strftime('%s', 'now'),?,?)
        `, fId, size, filename, content_type, sha1sum, uploadKey)

        if err == nil {
            break
        }
        fId = id.GenFileId()
    }

    if err != nil {
        tx.Rollback()
    }

    return fId, tx, err
}

func (db *DB) DeleteFile(id string) (bool, error) {
    res, err := db.ctx.Exec(`
        DELETE FROM files 
        WHERE id = ?;
    `, id)

    if err != nil {
        return false, err
    }
    n, err := res.RowsAffected()
    return n != 0, nil
}
