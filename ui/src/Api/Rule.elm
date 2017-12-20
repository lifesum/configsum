module Api.Rule exposing (activateRule, deactivateRule, getRule, listRules)

import Http
import Json.Decode as Decode
import Data.Rule exposing (Rule, decoder)


activateRule : String -> Http.Request String
activateRule id =
    Http.request
        { body = Http.emptyBody
        , expect = Http.expectStringResponse (returnId id)
        , headers = []
        , method = "PUT"
        , timeout = Nothing
        , url = "api/rules/" ++ id ++ "/activate"
        , withCredentials = False
        }


deactivateRule : String -> Http.Request String
deactivateRule id =
    Http.request
        { body = Http.emptyBody
        , expect = Http.expectStringResponse (returnId id)
        , headers = []
        , method = "PUT"
        , timeout = Nothing
        , url = "api/rules/" ++ id ++ "/deactivate"
        , withCredentials = False
        }


getRule : String -> Http.Request Rule
getRule id =
    Http.get ("api/rules/" ++ id) decoder


listRules : Http.Request (List Rule)
listRules =
    Http.get "api/rules/" (Decode.field "rules" (Decode.list decoder))



-- HELPER


returnId : String -> Http.Response String -> Result String String
returnId id response =
    if response.status.code == 204 then
        Ok id
    else
        Err response.status.message
