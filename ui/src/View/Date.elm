module View.Date exposing (pretty, short)

import Date exposing (Date)
import Time exposing (Time)


pretty : Time -> Date -> String
pretty now date =
    let
        day =
            24 * Time.hour

        diff =
            (now - (Date.toTime date))

        diffDays =
            Basics.floor (diff / day)
    in
        if diff < Time.minute then
            "just now"
        else if diff < (2 * Time.minute) then
            "1 minute ago"
        else if diff < Time.hour then
            (toString (Basics.floor (diff / Time.minute))) ++ " minutes ago"
        else if diff < (2 * Time.hour) then
            "1 hour ago"
        else if diff < day then
            (toString (Basics.floor (diff / Time.hour))) ++ " hours ago"
        else if diffDays == 1 then
            "yesterday"
        else if diffDays < 7 then
            (toString diffDays) ++ " days ago"
        else
            (toString (Basics.ceiling ((toFloat diffDays) / 7))) ++ " weeks ago"


short : Date -> String
short date =
    let
        month =
            case Date.month date of
                Date.Jan ->
                    "01"

                Date.Feb ->
                    "02"

                Date.Mar ->
                    "03"

                Date.Apr ->
                    "04"

                Date.May ->
                    "05"

                Date.Jun ->
                    "06"

                Date.Jul ->
                    "07"

                Date.Aug ->
                    "08"

                Date.Sep ->
                    "09"

                Date.Oct ->
                    "10"

                Date.Nov ->
                    "11"

                Date.Dec ->
                    "12"

        dateStr =
            String.join "-"
                [ (toString (Date.year date))
                , month
                , (toString (Date.day date))
                ]

        timeStr =
            String.join ":"
                [ (toString (Date.hour date))
                , (toString (Date.minute date))
                , (toString (Date.second date))
                ]
    in
        dateStr ++ " " ++ timeStr
