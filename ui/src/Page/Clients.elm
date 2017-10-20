module Page.Clients exposing (view)

import Html exposing (Html, div, h1, text)
import Html.Attributes exposing (class)


view : Html msg
view =
    div [ class "page" ]
        [ h1 [] [ text "Clients" ]
        ]
