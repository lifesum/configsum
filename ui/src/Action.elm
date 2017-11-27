module Action exposing (Msg(..))

import Time exposing (Time)
import Route exposing (Route)
import Page.Clients as Clients
import Page.Configs as Configs
import Page.Errored exposing (PageLoadError)


type Msg
    = ClientsLoaded (Result PageLoadError Clients.Model)
    | ClientsMsg Clients.Msg
    | ConfigsLoaded (Result PageLoadError Configs.Model)
    | ConfigsMsg Configs.Msg
    | ConfigBaseLoaded (Result PageLoadError Configs.Model)
    | LoadPage (Maybe Route)
    | SetRoute Route
    | Tick Time
