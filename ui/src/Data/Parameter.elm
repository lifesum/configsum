module Data.Parameter exposing (Parameter(..), decoder, paramsEncoder)

import Json.Decode as Decode exposing (Decoder)
import Json.Encode as Encode


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


decoder : Decoder Parameter
decoder =
    Decode.oneOf
        [ boolParamDecoder
        , numberParamDecoder
        , stringParamDecoder
        ]


encoder : Parameter -> ( String, Encode.Value )
encoder param =
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
        [ ( "parameters", Encode.object (List.map encoder params) )
        ]
