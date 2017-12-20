module Api.Config exposing (addParameter, createBase, getBase, listBase)

import Http
import Json.Decode as Decode
import Data.Config exposing (Config, decoder, encoder)
import Data.Parameter exposing (Parameter(..), paramsEncoder)


addParameter : String -> List Parameter -> Http.Request Config
addParameter id params =
    Http.request
        { body = (paramsEncoder params |> Http.jsonBody)
        , expect = Http.expectJson decoder
        , headers = []
        , method = "PUT"
        , timeout = Nothing
        , url = "api/configs/base/" ++ id
        , withCredentials = False
        }


createBase : String -> String -> Http.Request Config
createBase clientId name =
    Http.post "api/configs/base/" (encoder clientId name |> Http.jsonBody) decoder


getBase : String -> Http.Request Config
getBase id =
    Http.get ("api/configs/base/" ++ id) decoder


listBase : Http.Request (List Config)
listBase =
    Http.get "api/configs/base" (Decode.field "base_configs" (Decode.list decoder))
