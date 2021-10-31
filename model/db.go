package model

import (
    "busmap.vn/librarycore/config"
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    "github.com/qor/media"
    "os"
)

var DBInstance *gorm.DB

func SetUpDbInstance(dbInst *gorm.DB) {
    // =========
    // METHOD 1
    // =========
    if dbInst == nil {
        panic("[DB][LIBRARY] Error: \"dbInst\" mustn't be null!")
    } else {
        DBInstance = dbInst
    }

    // Create necessary extension
    DBInstance.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm")
    DBInstance.Exec("CREATE EXTENSION IF NOT EXISTS unaccent")

    // Migrate tables
    SkipDBMigration := os.Getenv("SKIP_DB_MIGRATION")
    if SkipDBMigration == "1" {
        fmt.Println("[DB][LIBRARY] SKIP_DB_MIGRATION has activated")
    } else {
        fmt.Println("[DB][LIBRARY] Migrating library tables...")
        _autoMigrateTables()
        fmt.Println("[DB][LIBRARY] Migrating process has done!")
    }

    // Initialize configs for tables
    createNecessaryConfigsForTables := os.Getenv("CREATE_NECESSARY_CONFIGS_FOR_TABLES")
    if createNecessaryConfigsForTables == "1" {
        _initializeTableConfig()
    }
}

func Initialize() {
    // =========
    // METHOD 2
    // =========

    databaseName := "Postgres (Production)"

    // ---------------------

    // Connect to Database
    var err error
    DBInstance, err = gorm.Open("postgres", config.Config.PostgresConnectionString)
    if err != nil {
        panic("[DB] Error: Open database fail!")
    }

    if DBInstance != nil && DBInstance.DB().Ping() == nil {
        fmt.Println("Yay! Database " + databaseName + " has connected successfully!")
        DBInstance.LogMode(config.Config.PostgresLogMode)
    } else {
        errMsg := fmt.Sprintf("[DB][Error] Database (%v) connected fail!\n", databaseName)
        panic(errMsg)
    }

    // ---------------------

    // Create necessary extension
    DBInstance.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm")
    DBInstance.Exec("CREATE EXTENSION IF NOT EXISTS unaccent")

    // ---------------------

    SkipDBMigration := os.Getenv("SKIP_DB_MIGRATION")
    if SkipDBMigration == "1" {
        fmt.Println("[DB] SKIP_DB_MIGRATION has activated")
    } else {
        fmt.Println("[DB] Migrating system & library tables...")
        _autoMigrateTables()
        fmt.Println("[DB] Migrating process has done!")
    }

    // ---------------------

    media.RegisterCallbacks(DBInstance)

    // ---------------------

    if config.Config.CreateNecessaryConfigsForTables == "1" {
        _initializeTableConfig()
    }
}

func _autoMigrateTables() {
    _ = DBInstance.AutoMigrate(
        // TODO: Migrate other library tables here !!!
        &Classroom{},
    )
}

func _initializeTableConfig() {
    // Unaccent function in database
    // Bởi vì không thể tạo index "unaccent(lower(name))" nên phải
    // tạo function với IMMUTABLE. Từ đây, tạo ra index.
    // Nguồn tham khảo: https://stackoverflow.com/questions/11005036/does-postgresql-support-accent-insensitive-collations/11007216#11007216
    DBInstance.Exec(`
    CREATE OR REPLACE FUNCTION f_unaccent(text)
    RETURNS TEXT AS $func$
    DECLARE input_string text := LOWER($1);
    BEGIN
        input_string := translate(input_string, 'áàãạảAÁÀÃẠẢăắằẵặẳĂẮẰẴẶẲâầấẫậẩÂẤẦẪẬẨ', 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa');
        input_string := translate(input_string, 'éèẽẹẻEÉÈẼẸẺêếềễệểÊẾỀỄỆỂ', 'eeeeeeeeeeeeeeeeeeeeeeee');
        input_string := translate(input_string, 'íìĩịỉIÍÌĨỊỈ', 'iiiiiiiiiii');
        input_string := translate(input_string, 'óòõọỏOÓÒÕỌỎôốồỗộổÔỐỒỖỘỔơớờỡợởƠỚỜỠỢỞ', 'ooooooooooooooooooooooooooooooooooo');
        input_string := translate(input_string, 'úùũụủUÚÙŨỤỦưứừữựửƯỨỪỮỰỬ', 'uuuuuuuuuuuuuuuuuuuuuuu');
        input_string := translate(input_string, 'ýỳỹỵỷYÝỲỸỴỶ', 'yyyyyyyyyyy');
        input_string := translate(input_string, 'dđĐD', 'dddd');
        return input_string;
    END;

    $func$ LANGUAGE plpgsql IMMUTABLE`)
}
