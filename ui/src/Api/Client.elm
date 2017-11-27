module Api.Client exposing (create, list)

import Http
import Json.Decode as Decode
import Data.Client exposing (Client, decoder, encoder)


create : String -> Http.Request Client
create name =
    Http.post "api/clients/" (encoder name |> Http.jsonBody) decoder


list : Http.Request (List Client)
list =
    Http.get "api/clients/" (Decode.field "clients" (Decode.list decoder))
