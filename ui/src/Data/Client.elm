module Data.Client exposing (Client, decoder, encoder)

import Date exposing (Date)
import Json.Decode as Decode exposing (Decoder, andThen, fail, succeed)
import Json.Encode as Encode


type alias Client =
    { createdAt : Date
    , deleted : Bool
    , id : String
    , name : String
    , token : String
    }


decoder : Decoder Client
decoder =
    Decode.map5 Client
        (Decode.field "created_at" date)
        (Decode.field "deleted" Decode.bool)
        (Decode.field "id" Decode.string)
        (Decode.field "name" Decode.string)
        (Decode.field "token" Decode.string)


encoder : String -> Encode.Value
encoder name =
    Encode.object
        [ ( "name", Encode.string name )
        ]



-- HELPER


date : Decoder Date
date =
    let
        convert : String -> Decoder Date
        convert raw =
            case Date.fromString raw of
                Ok date ->
                    succeed date

                Err error ->
                    fail error
    in
        Decode.string |> andThen convert
