module Page.Clients exposing (Model, Msg, init, update, view)

import Html exposing (Html, a, button, div, form, h1, input, span, table, tbody, td, text, th, thead, tr)
import Html.Attributes exposing (class, placeholder, type_, value)
import Html.Events exposing (onClick, onInput, onSubmit)
import Http
import Task exposing (Task)
import Api.Client as Api
import Data.Client exposing (Client)
import Page.Errored exposing (PageLoadError, pageLoadError)
import View.Error


-- MDOEL


type alias Model =
    { clients : List Client
    , error : Maybe Http.Error
    , formName : String
    , showCreate : Bool
    }


init : Task PageLoadError Model
init =
    Api.list
        |> Http.toTask
        |> Task.map (\clients -> Model clients Nothing "" False)
        |> Task.mapError (\err -> pageLoadError "Clients" err)



-- UPDATE


type Msg
    = FormSubmit
    | FormSubmitted (Result Http.Error Client)
    | ToggleCreate
    | UpdateName String


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        FormSubmit ->
            ( { model | error = Nothing }, Api.create model.formName |> Http.send FormSubmitted )

        FormSubmitted (Err error) ->
            ( { model | error = Just error }, Cmd.none )

        FormSubmitted (Ok client) ->
            ( Model (client :: model.clients) Nothing "" False, Cmd.none )

        ToggleCreate ->
            ( { model | error = Nothing, showCreate = not model.showCreate }, Cmd.none )

        UpdateName name ->
            ( { model | error = Nothing, formName = name }, Cmd.none )



-- VIEW


view : Model -> Html Msg
view model =
    div []
        [ h1 [] [ text "Clients" ]
        , View.Error.view model.error
        , viewCreate model.formName model.showCreate
        , viewList model.clients
        ]


viewCreate : String -> Bool -> Html Msg
viewCreate name showForm =
    let
        content =
            if showForm then
                form [ onSubmit FormSubmit ]
                    [ input
                        [ onInput UpdateName
                        , placeholder "Name"
                        , type_ "text"
                        , value name
                        ]
                        []
                    , button [] [ text "create" ]
                    ]
            else
                a
                    [ class "action"
                    , onClick ToggleCreate
                    ]
                    [ span [ class "nc-icon nc-circle-add" ] []
                    , span [] [ text "create client" ]
                    ]
    in
        div [ class "create" ] [ content ]


viewList : List Client -> Html Msg
viewList clients =
    table []
        [ thead []
            [ tr []
                [ th [] [ text "name" ]
                , th [] [ text "token" ]
                ]
            ]
        , tbody [] (List.map viewItem clients)
        ]


viewItem : Client -> Html Msg
viewItem client =
    tr []
        [ td [] [ text client.name ]
        , td [] [ text client.token ]
        ]
