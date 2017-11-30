module Main exposing (main)

import Html exposing (Html, div, footer, h1, text)
import Html.Attributes exposing (class)
import Navigation
import Task
import Time exposing (Time)
import Action exposing (Msg(..))
import Page.Blank as Blank
import Page.Clients as Clients
import Page.Configs as Configs
import Page.Errored as Errored exposing (PageLoadError)
import Page.Rules as Rules
import Route exposing (Route, navigate)
import View.Page as Page


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
    { now : Time
    }


type alias Model =
    { pageState : PageState
    , route : Route
    , now : Time
    }


type Page
    = Blank String
    | Clients Clients.Model
    | Configs Configs.Model
    | Rules Rules.Model
    | Errored PageLoadError
    | NotFound


type PageState
    = Loaded Page
    | TransitioningFrom Page


init : Flags -> Navigation.Location -> ( Model, Cmd Msg )
init { now } location =
    let
        route =
            case Route.fromLocation location of
                Nothing ->
                    Route.NotFound

                Just route ->
                    route

        model =
            Model (Loaded (Blank "Loading")) route now
    in
        setRoute (Route.fromLocation location) model


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.batch
        [ Time.every Time.second Tick
        ]



-- UPDATE --


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    let
        page =
            getPage model.pageState

        toPage toModel toMsg subUpdate subMsg subModel =
            let
                ( newModel, newCmd ) =
                    subUpdate subMsg subModel
            in
                ( { model | pageState = Loaded (toModel newModel) }, Cmd.map toMsg newCmd )
    in
        case ( (Debug.log "MSG" msg), page ) of
            ( ClientsLoaded (Err error), _ ) ->
                ( { model | pageState = Loaded (Errored error) }, Cmd.none )

            ( ClientsLoaded (Ok subModel), _ ) ->
                ( { model | pageState = Loaded (Clients subModel) }, Cmd.none )

            ( ClientsMsg subMsg, Clients subModel ) ->
                toPage Clients ClientsMsg Clients.update subMsg subModel

            ( ConfigsLoaded (Err error), _ ) ->
                ( { model | pageState = Loaded (Errored error) }, Cmd.none )

            ( ConfigsLoaded (Ok subModel), _ ) ->
                ( { model | pageState = Loaded (Configs subModel) }, Cmd.none )

            ( ConfigsMsg subMsg, Configs subModel ) ->
                toPage Configs ConfigsMsg Configs.update subMsg subModel

            ( ConfigBaseLoaded (Err error), _ ) ->
                ( { model | pageState = Loaded (Errored error) }, Cmd.none )

            ( ConfigBaseLoaded (Ok subModel), _ ) ->
                ( { model | pageState = Loaded (Configs subModel) }, Cmd.none )

            ( LoadPage maybeRoute, _ ) ->
                let
                    route =
                        case maybeRoute of
                            Nothing ->
                                Route.Clients

                            Just route ->
                                route
                in
                    setRoute maybeRoute { model | route = route }

            ( RulesLoaded (Err error), _ ) ->
                ( { model | pageState = Loaded (Errored error) }, Cmd.none )

            ( RulesLoaded (Ok subModel), _ ) ->
                ( { model | pageState = Loaded (Rules subModel) }, Cmd.none )

            ( RuleLoaded (Err error), _ ) ->
                ( { model | pageState = Loaded (Errored error) }, Cmd.none )

            ( RuleLoaded (Ok subModel), _ ) ->
                ( { model | pageState = Loaded (Rules subModel) }, Cmd.none )

            ( RulesMsg subMsg, Rules subModel ) ->
                toPage Rules RulesMsg Rules.update subMsg subModel

            ( SetRoute route, _ ) ->
                ( model, navigate route )

            ( Tick now, loadedPage ) ->
                let
                    newPage =
                        case loadedPage of
                            Configs configsModel ->
                                Configs ({ configsModel | now = now })

                            _ ->
                                loadedPage
                in
                    ( { model | pageState = (Loaded newPage), now = now }, Cmd.none )

            ( _, NotFound ) ->
                ( model, Cmd.none )

            ( _, _ ) ->
                ( model, Cmd.none )


setRoute : Maybe Route.Route -> Model -> ( Model, Cmd Msg )
setRoute maybeRoute model =
    case maybeRoute of
        Nothing ->
            ( { model | pageState = Loaded NotFound }, Cmd.none )

        Just Route.Clients ->
            ( { model | pageState = TransitioningFrom <| getPage model.pageState }
            , Task.attempt ClientsLoaded Clients.init
            )

        Just Route.Configs ->
            ( model, navigate Route.ConfigsBase )

        Just Route.ConfigsBase ->
            ( { model | pageState = TransitioningFrom <| getPage model.pageState }
            , Task.attempt ConfigsLoaded (Configs.init model.now)
            )

        Just (Route.ConfigBase id) ->
            ( { model | pageState = TransitioningFrom <| getPage model.pageState }
            , Task.attempt ConfigBaseLoaded <| Configs.initBase model.now id
            )

        Just Route.NotFound ->
            ( { model | pageState = Loaded NotFound }, Cmd.none )

        Just Route.Rules ->
            ( { model | pageState = TransitioningFrom <| getPage model.pageState }
            , Task.attempt RulesLoaded <| Rules.initList model.now
            )

        Just (Route.Rule id) ->
            ( { model | pageState = TransitioningFrom <| getPage model.pageState }
            , Task.attempt RuleLoaded <| Rules.initRule model.now id
            )



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
            Blank name ->
                Blank.view name
                    |> frame name

            Clients subModel ->
                Clients.view subModel
                    |> Html.map ClientsMsg
                    |> frame "clients"

            Configs subModel ->
                Configs.view subModel
                    |> Html.map ConfigsMsg
                    |> frame "configs"

            Errored error ->
                Errored.view error
                    |> frame "errored"

            NotFound ->
                Blank.view "Not Found"
                    |> frame "not-found"

            Rules subModel ->
                Rules.view subModel
                    |> Html.map RulesMsg
                    |> frame "rules"



-- HELPER


getPage : PageState -> Page
getPage state =
    case state of
        Loaded page ->
            page

        TransitioningFrom page ->
            page
