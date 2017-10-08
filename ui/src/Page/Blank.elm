module Page.Blank exposing (view)

import Html exposing (Html, div, h1, text)
import Html.Attributes exposing (class)


view : String -> Html msg
view name =
    div [ class "page" ]
        [ h1 [] [ text name ]
        ]
