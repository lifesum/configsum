module View.Error exposing (view)

import Html exposing (Html, div, text)
import Html.Attributes exposing (class)
import Http
import Json.Decode as Decode
import Data.Error as Error


view : Maybe Http.Error -> Html msg
view error =
    case error of
        Just error ->
            div [ class "error" ] [ text (errorMessage error) ]

        Nothing ->
            div [] []



-- HELPER


errorMessage : Http.Error -> String
errorMessage httpError =
    case httpError of
        Http.BadStatus response ->
            case (Decode.decodeString Error.decoder response.body) of
                Err decodeError ->
                    decodeError

                Ok error ->
                    error.reason

        _ ->
            "Something went wrong!"
