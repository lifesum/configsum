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
                [ viewNavLink (activeRoute route) Route.Clients [ text "clients" ]
                , viewNavLink (activeRoute route) Route.Configs [ text "configs" ]
                , viewNavLink (activeRoute route) Route.Rules [ text "rules" ]
                ]
            ]
        ]


viewFooter : Html Msg
viewFooter =
    footer [] []


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


activeRoute : Route -> Route
activeRoute route =
    case route of
        Route.ConfigsBase ->
            Route.Configs

        Route.ConfigBase _ ->
            Route.Configs

        _ ->
            route


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
