package constants

const (
    GenderUnknown = 0
    GenderMale    = 1
    GenderFemale  = 2
    // ------------------------------
    OsAll     = "all"
    OsAndroid = "android"
    OsIos     = "ios"
    // ------------------------------
    DefaultLocale     = "vi"
    LocaleAll         = "all"
    LocaleVietnamese  = "vi"
    LocaleEnglish     = "en"
    LocationName      = "Asia/Ho_Chi_Minh"
    LocationGMT       = 7
    NoonTimeBy24Hour  = 12 // By GMT+7 (VN), noon time will be 12:00 PM
    TimeZeroValue     = "0001-01-01 00:00:00+00"
    TimeZeroValueUnix = -62135622000
    InMorning         = 1
    InAfternoon       = 2
    // ------------------------------
    MaxHumanAge = 100
    // ------------------------------
    UPLOAD_FILE_FROM_APP    = "/system/upload/library"
    EXPORT_USER_FOLDER_PATH = "/system/export/library"
    // ------------------------------
    Action_OpenBorrowingDetail    = "LIBRARY_OPEN_BORROWING_DETAIL"
    Action_OpenRegistrationDetail = "LIBRARY_OPEN_BORROWING_REGISTRATION_DETAIL"
    Action_OpenBookHeadDetail     = "LIBRARY_OPEN_BOOK_HEAD_DETAIL"
    Action_OpenCardList           = "LIBRARY_OPEN_CARD_LIST"
    Action_OpenUrl                = "LIBRARY_OPEN_URL"
    Action_OpenUrlInApp           = "LIBRARY_OPEN_URL_INAPP"
)

// =============================================================
// =============================================================
// =============================================================

func GetGenderArray() []uint {
    return []uint{
        GenderUnknown,
        GenderMale,
        GenderFemale,
    }
}

func GetOsArray() []string {
    return []string{
        OsAll,
        OsAndroid,
        OsIos,
    }
}

func GetAllLocaleArray() []string {
    return []string{
        LocaleAll,
        LocaleVietnamese,
        LocaleEnglish,
    }
}

func GetValidLocaleArray() []string {
    return []string{
        LocaleVietnamese,
        LocaleEnglish,
    }
}

func GetContentByLanguage(
    language string,
    contentVi string,
    contentEn string,
) string {
    switch language {
    case LocaleVietnamese:
        return contentVi
    case LocaleEnglish:
        return contentEn
    default:
        return contentVi
    }
}
