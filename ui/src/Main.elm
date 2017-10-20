module Main exposing (main)

import Html exposing (Html, div, footer, h1, text)
import Html.Attributes exposing (class)
import Navigation
import Action exposing (Msg(..))
import Page.Clients as Clients
import Page.Blank as Blank
import Route exposing (Route, navigate)
import Views.Page as Page


-- MAIN --


main : Program Flags Model Msg
main =
    Navigation.programWithFlags (Route.fromLocation >> LoadPage)
        { init = init
        , subscriptions = subscriptions
        , update = update
        , view = view
        }



-- MODEL --


type alias Flags =
    {}


type alias Model =
    { pageState : PageState
    , route : Route
    }


type Page
    = Blank String
    | Clients
    | NotFound


type PageState
    = Loaded Page
    | TransitioningFrom Page


init : Flags -> Navigation.Location -> ( Model, Cmd Msg )
init _ location =
    let
        route =
            case Route.fromLocation location of
                Nothing ->
                    Route.NotFound

                Just route ->
                    route

        model =
            { pageState = Loaded (Blank "Loading"), route = route }
    in
        setRoute (Route.fromLocation location) model


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.batch []



-- UPDATE --


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case (Debug.log "MSG" msg) of
        LoadPage maybeRoute ->
            let
                route =
                    case maybeRoute of
                        Nothing ->
                            Route.Clients

                        Just route ->
                            route
            in
                setRoute maybeRoute { model | route = route }

        SetRoute route ->
            ( model, navigate route )

        Tick _ ->
            ( model, Cmd.none )


setRoute : Maybe Route.Route -> Model -> ( Model, Cmd Msg )
setRoute maybeRoute model =
    case maybeRoute of
        Nothing ->
            ( { model | pageState = Loaded NotFound }, Cmd.none )

        Just Route.Clients ->
            ( { model | pageState = Loaded Clients }, Cmd.none )

        Just Route.Configs ->
            ( { model | pageState = Loaded (Blank "Configs") }, Cmd.none )

        Just Route.NotFound ->
            ( { model | pageState = Loaded NotFound }, Cmd.none )

        Just Route.Rules ->
            ( { model | pageState = Loaded (Blank "Rules") }, Cmd.none )



-- VIEW --


view : Model -> Html Msg
view model =
    let
        content =
            case model.pageState of
                Loaded page ->
                    viewPage False page model.route

                TransitioningFrom page ->
                    viewPage True page model.route
    in
        div []
            [ content
            , footer []
                [ div [ class "debug" ] [ text (toString model) ]
                ]
            ]


viewPage : Bool -> Page -> Route -> Html Msg
viewPage isLoading page route =
    let
        frame =
            Page.frame isLoading route
    in
        case page of
            Clients ->
                Clients.view
                    |> frame

            Blank name ->
                Blank.view name
                    |> frame

            NotFound ->
                Blank.view "Not Found"
                    |> frame
