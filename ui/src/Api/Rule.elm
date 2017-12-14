module Api.Rule exposing (getRule)

import Http
import Data.Rule exposing (Rule, decoder)


getRule : String -> Http.Request Rule
getRule id =
    Http.get ("api/rules/" ++ id) decoder
