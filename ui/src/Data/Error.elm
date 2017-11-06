module Data.Error exposing (Error, decoder)

import Json.Decode as Decode exposing (Decoder)


type alias Error =
    { reason : String
    }


decoder : Decoder Error
decoder =
    Decode.map Error
        (Decode.field "reason" Decode.string)
