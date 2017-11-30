module Data.Config exposing (Config, decoder, encoder)

import Date exposing (Date)
import Json.Decode as Decode exposing (Decoder, andThen, fail, succeed)
import Json.Encode as Encode
import Data.Parameter exposing (Parameter(..))


type alias Config =
    { clientId : String
    , id : String
    , name : String
    , parameters : List Parameter
    , createdAt : Date
    , updatedAt : Date
    }


decoder : Decoder Config
decoder =
    Decode.map6 Config
        (Decode.field "client_id" Decode.string)
        (Decode.field "id" Decode.string)
        (Decode.field "name" Decode.string)
        (Decode.field "parameters" (Decode.list Data.Parameter.decoder))
        (Decode.field "created_at" date)
        (Decode.field "updated_at" date)


encoder : String -> String -> Encode.Value
encoder clientId name =
    Encode.object
        [ ( "client_id", Encode.string clientId )
        , ( "name", Encode.string name )
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
