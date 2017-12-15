module Api.Rule exposing (getRule, listRules)

import Http
import Json.Decode as Decode
import Data.Rule exposing (Rule, decoder)


getRule : String -> Http.Request Rule
getRule id =
    Http.get ("api/rules/" ++ id) decoder


listRules : Http.Request (List Rule)
listRules =
    Http.get "api/rules/" (Decode.field "rules" (Decode.list decoder))
