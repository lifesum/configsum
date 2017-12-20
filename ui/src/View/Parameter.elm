module View.Parameter exposing (name, typeClass, viewTable, viewTableItem)

import Html exposing (Html, h2, div, input, label, section, span, table, tbody, td, text, th, thead, tr)
import Html.Attributes exposing (checked, class, classList, disabled, for, id, type_)
import Data.Parameter exposing (Parameter(..))


-- VIEW


viewTable : List (Html msg) -> List Parameter -> Html msg
viewTable action params =
    section [ class "parameters" ]
        [ h2 [] [ text "parameters" ]
        , table []
            [ thead []
                [ tr []
                    [ th [ class "name" ] [ text "name" ]
                    , th [ class "type" ] [ text "type" ]
                    , th [] [ text "value" ]
                    ]
                ]
            , tbody [] <| List.append (List.map viewTableItem params) action
            ]
        ]


viewTableItem : Parameter -> Html msg
viewTableItem parameter =
    tr []
        [ td [] [ text <| name parameter ]
        , td
            [ classList [ ( "type", True ), ( (typeClass parameter), True ) ]
            ]
            [ text <| typeClass parameter
            ]
        , td
            [ class ("value " ++ (typeClass parameter))
            ]
            [ viewTableItemValue parameter
            ]
        ]


viewTableItemValue : Parameter -> Html msg
viewTableItemValue parameter =
    case parameter of
        BoolParameter name value ->
            div []
                [ input
                    [ checked value
                    , disabled True
                    , id ("param-bool-" ++ name)
                    , type_ "checkbox"
                    ]
                    []
                , label [ for ("param-bool-" ++ name) ] []
                ]

        NumberParameter _ value ->
            span [] [ text (toString value) ]

        NumbersParameter _ values ->
            div [] (List.map (\v -> span [] [ text (toString v) ]) values)

        StringParameter _ value ->
            span [] [ text value ]

        StringsParameter _ values ->
            div [] (List.map (\v -> span [] [ text v ]) values)



-- HELPER


name : Parameter -> String
name parameter =
    case parameter of
        BoolParameter name _ ->
            name

        NumberParameter name _ ->
            name

        NumbersParameter name _ ->
            name

        StringParameter name _ ->
            name

        StringsParameter name _ ->
            name


typeClass : Parameter -> String
typeClass parameter =
    case parameter of
        BoolParameter _ _ ->
            "bool"

        NumberParameter _ _ ->
            "number"

        NumbersParameter _ _ ->
            "numbers"

        StringParameter _ _ ->
            "string"

        StringsParameter _ _ ->
            "strings"
