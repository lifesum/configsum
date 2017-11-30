module Page.Rules exposing (Model, Msg, init, update, view)

import Html exposing (Html, div, h1, table, text)
import Http
import Task exposing (Task)
import Data.Rule exposing (Rule)
import Page.Errored exposing (PageLoadError)
import View.Error


-- MODEL


type alias Model =
    { error : Maybe Http.Error
    , rule : Maybe Rule
    }


init : Task PageLoadError Model
init =
    Task.succeed <| Model Nothing Nothing



-- UPDATE


type Msg
    = Noop


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    ( model, Cmd.none )



-- VIEW


view : Model -> Html Msg
view model =
    case model.rule of
        Just rule ->
            viewRule model

        Nothing ->
            viewList model


viewList : Model -> Html Msg
viewList model =
    div []
        [ h1 [] [ text "Rules" ]
        , View.Error.view model.error
        , table []
            []
        ]


viewRule : Model -> Html Msg
viewRule model =
    div [] []
