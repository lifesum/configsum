module Page.Errored exposing (PageLoadError, pageLoadError, view)

import Html exposing (Html, div, h1, h2, text)
import Html.Attributes exposing (class)
import Http


type PageLoadError
    = PageLoadError String Http.Error


view : PageLoadError -> Html msg
view (PageLoadError page err) =
    div []
        [ unpackError err
        ]


pageLoadError : String -> Http.Error -> PageLoadError
pageLoadError page err =
    PageLoadError page err


unpackError : Http.Error -> Html msg
unpackError err =
    case err of
        Http.BadPayload debug res ->
            div [ class "error" ]
                [ h2 [] [ text debug ]
                , div [] [ text (toString res) ]
                ]

        Http.BadStatus res ->
            div [ class "error" ]
                [ h2 [] [ text "Entity not found" ]
                , div [] [ text (toString res) ]
                ]

        _ ->
            div [ class "error" ] [ text "Unhandled error" ]
