module View.Page exposing (frame)

import Html exposing (Attribute, Html, a, div, header, footer, li, main_, nav, span, text, ul)
import Html.Attributes exposing (class, classList)
import Html.Events exposing (defaultOptions, onWithOptions)
import Json.Decode exposing (Decoder, andThen, fail, map2, succeed)
import Action exposing (Msg(..))
import Route exposing (Route)


frame : Bool -> Route -> String -> Html Msg -> Html Msg
frame isLoading route page content =
    div []
        [ viewHeader route
        , main_ [ class ("page " ++ page) ]
            [ content
            ]
        ]


viewHeader : Route -> Html Msg
viewHeader route =
    header []
        [ div [ class "logo" ]
            [ span [ class "lead" ] [ text "C" ]
            , span [] [ text "onfigsum" ]
            ]
        , nav []
            [ ul []
                [ viewNavLink (Route.active route) Route.Clients [ text "clients" ]
                , viewNavLink (Route.active route) Route.Configs [ text "configs" ]
                , viewNavLink (Route.active route) Route.Rules [ text "rules" ]
                ]
            ]
        ]


viewNavLink : Route -> Route -> List (Html Msg) -> Html Msg
viewNavLink activeRoute route content =
    li []
        [ a
            [ classList [ ( "active", activeRoute == route ) ]
            , Route.href route
            , onClickRoute route
            ]
            content
        ]



-- ATTRIBUTES


onClickRoute : Route -> Attribute Msg
onClickRoute route =
    onWithOptions
        "click"
        { defaultOptions | preventDefault = True }
        (map2
            invertedOr
            (Json.Decode.field "ctrlKey" Json.Decode.bool)
            (Json.Decode.field "metaKey" Json.Decode.bool)
            |> andThen (maybePreventDefault (SetRoute route))
        )



-- HELPER


maybePreventDefault : Msg -> Bool -> Decoder Msg
maybePreventDefault msg isPrevent =
    case isPrevent of
        True ->
            succeed msg

        False ->
            fail "Normal Link"


invertedOr : Bool -> Bool -> Bool
invertedOr x y =
    not (x || y)
