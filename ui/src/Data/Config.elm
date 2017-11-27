module Data.Config exposing (Config, Parameter(..), decoder, encoder, paramsEncoder)

import Date exposing (Date)
import Json.Decode as Decode exposing (Decoder, andThen, fail, succeed)
import Json.Encode as Encode


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
        (Decode.field "parameters" (Decode.list paramDecoder))
        (Decode.field "created_at" date)
        (Decode.field "updated_at" date)


encoder : String -> String -> Encode.Value
encoder clientId name =
    Encode.object
        [ ( "client_id", Encode.string clientId )
        , ( "name", Encode.string name )
        ]


type Parameter
    = BoolParameter String Bool
    | NumberParameter String Int
    | NumbersParameter String (List Int)
    | StringParameter String String
    | StringsParameter String (List String)


boolParamDecoder : Decoder Parameter
boolParamDecoder =
    Decode.map2 BoolParameter
        (Decode.field "name" Decode.string)
        (Decode.field "value" Decode.bool)


numberParamDecoder : Decoder Parameter
numberParamDecoder =
    Decode.map2 NumberParameter
        (Decode.field "name" Decode.string)
        (Decode.field "value" Decode.int)


stringParamDecoder : Decoder Parameter
stringParamDecoder =
    Decode.map2 StringParameter
        (Decode.field "name" Decode.string)
        (Decode.field "value" Decode.string)


paramDecoder : Decoder Parameter
paramDecoder =
    Decode.oneOf
        [ boolParamDecoder
        , numberParamDecoder
        , stringParamDecoder
        ]


paramEncode : Parameter -> ( String, Encode.Value )
paramEncode param =
    case param of
        BoolParameter name value ->
            ( name, Encode.bool value )

        NumberParameter name value ->
            ( name, Encode.int value )

        StringParameter name value ->
            ( name, Encode.string value )

        _ ->
            ( "", Encode.string "unknown type" )


paramsEncoder : List Parameter -> Encode.Value
paramsEncoder params =
    Encode.object
        [ ( "parameters", Encode.object (List.map paramEncode params) )
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
