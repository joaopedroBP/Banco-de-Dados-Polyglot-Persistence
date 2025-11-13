#!/usr/bin/env bash

set -e

# mostp = most playem
accepted_subcommands=("user" "artist" "album" "track" "playlist" "genre" "history" "like" "mostp" "time")

# add artist to DB
function handleadd_artist() {
  if [[ -z "$4" ]]; then
    echo "Usage: $0 add artist <artist_name> <description>"
    exit 1
  fi
  artist_name="$3"
  description="$4"
  # prints artist id
  ./s2postgre.sh add artist "$artist_name" "$description"
  exit 0
}

# add album to DB
function handleadd_album() {
  if [[ -z "$6" ]]; then
    echo "Usage: $0 add album <album_name> <artist_id> <release_date> <genre_id>"
    exit 1
  fi
  album_name="$3"
  artist_id="$4"
  release_year="$5"
  genre_id="$6"
  # prints album id
  ./s2postgre.sh add album "$album_name" "$artist_id" "$release_year" "$genre_id"
  exit 0
}

# add track to DB
function handleadd_track() {
  if [[ -z "$6" ]]; then
    echo "Usage: $0 add track <track_name> <album_id> <duration> <album_num>"
    exit 1
  fi
  track_name="$3"
  album_id="$4"
  duration="$5"
  album_num="$6"
  # prints track id
  ./s2postgre.sh add track "$track_name" "$album_id" "$duration" "$album_num"
  exit 0 
}

# add playlist to DB
function handleadd_playlist() {
  if [[ -z "$5" ]]; then
    echo "Usage: $0 add playlist <user_id> <playlist_name> <is_public>"
    echo "Note: is_public should be 'true' or 'false'"
    exit 1
  fi
  user_id="$3"
  playlist_name="$4"
  description="$5"
  # prints playlist id
  ./../S2_Mongo/S2_Mongo add playlist "$user_id" "$playlist_name" "$is_public" 
  exit 0
}

# add genre to DB
function handleadd_genre() {
  if [[ -z "$3" ]]; then
    echo "Usage: $0 add genre <genre_name>"
    exit 1
  fi
  genre_name="$3"
  # prints genre id
  ./s2postgre.sh add genre "$genre_name"
  exit 0
}

# add song to user's listening history
function handleadd_history() {
  if [[ -z "$4" ]]; then
    echo "Usage: $0 add history <user_id> <track_id>"
    exit 1
  fi
  user_id="$3"
  if [[ -z "$(./s2postgre.sh list user | grep "User: $user_id")" ]]; then
    echo "User not found"
  fi
  track_id="$4"
  if [[ -z "$(./s2postgre.sh get track "$track_id" | grep "ID: $track_id")" ]]; then
    echo "Track $track_id does not exist."
    exit 1
  fi
  # prints confirmation
  ./../S2_Scylla/S2_Scylla add "history" "$user_id" "$track_id" 
  exit 0
}

# add like to user's liked songs
function handleadd_like() {
  if [[ -z "$4" ]]; then
    echo "Usage: $0 add like <user_id> <track_id>"
    exit 1
  fi
  user_id="$3"
  track_id="$4"
  # check user
  if [[ -z "$(./s2postgre.sh list user | grep "User: $user_id")" ]]; then
    echo "User $user_id does not exist."
    exit 1
  fi
  # check track
  if [[ -z "$(./s2postgre.sh get track "$track_id" | grep "ID: $track_id")" ]]; then
    echo "Track $track_id does not exist."
    exit 1
  fi
  # prints confirmation
  ./../S2_Mongo/S2_Mongo add track "$user_id" likes "$track_id"
  exit 0
}

# create likes playlist for new user
function handleadd_likeslist() {
  user_id="$1"
  # check user exists
  if [[ -z "$(./s2postgre.sh list user | grep "User: $user_id")" ]]; then
    echo "User $user_id does not exist."
    exit 1                                                                NAMES
5d4f1ad0c770   scylladb/scylla:5.2                           "/docker-entrypoint.…"   13 hours ago   Up 2 hours   22/tcp, 7000-7001/tcp, 9160/tcp, 918
  fi
  # prints playlist creation confirmation
  ./../S2_Mongo/S2_Mongo add playlist "$user_id" "likes" "false"
  exit 0
}

# add time for a track for a user
function handleadd_time() {
  if [[ -z "$5" ]]; then
    echo "Usage: $0 add time <user_id> <track_id> <time_seconds>"
    exit 1
  fi
  user_id="$3"
  # check user
  if [[ -z "$(./s2postgre.sh list user | grep "User: $user_id")" ]]; then
    echo "User $user_id does not exist."
    exit 1
  fi
  track_id="$4"
  # check track
  if [[ -z "$(./s2postgre.sh get track "$track_id" | grep "ID: $track_id")" ]]; then
    echo "Track $track_id does not exist."
    exit 1
  fi
  time_seconds="$5"
  # prints confirmation
  ./../S2_Scylla/S2_Scylla add "time" "$user_id" "$track_id" "$time_seconds"
  handleadd_history "$@"
  exit 0
}

# add user to DB
function handleadd_user() {
  if [[ -z "$5" ]]; then
    echo "Usage: $0 add user <username> <email> <password>"
    exit 1
  fi
  username="$3"
  email="$4"
  password="$5"
  ./s2postgre.sh add user "$3" "$4" "$5"
  # check if user exists
  if [[ -z "$(./s2postgre.sh list user | grep "User: $username")" ]]; then
    echo "Failed to add user."
    exit 1
  fi
  handleadd_likeslist "$username"
  exit 0
}

# remove user from DB
function handlerm_user() {
  if [[ -z "$3" ]]; then
    echo "Usage: $0 "rm" user <user_id>"
    exit 1
  fi
  user_id="$3"
  ./s2postgre.sh "rm" user "$user_id"
  exit 0
}

# remove artist from DB
function handlerm_artist() {
  if [[ -z "$3" ]]; then
    echo "Usage: $0 "rm" artist <artist_id>"
    exit 1
  fi
  artist_id="$3"
  ./s2postgre.sh "rm" artist "$artist_id"
  exit 0
}

# remove album from DB
function handlerm_album() {
  if [[ -z "$3" ]]; then
    echo "Usage: $0 "rm" album <album_id>"
    exit 1
  fi
  album_id="$3"
  ./s2postgre.sh "rm" album "$album_id"
  exit 0
}

# remove track from DB
function handlerm_track() {
  if [[ -z "$3" ]]; then
    echo "Usage: $0 "rm" track <track_id>"
    exit 1
  fi
  track_id="$3"
  ./s2postgre.sh "rm" track "$track_id"
  exit 0
}

# remove playlist from DB
function handlerm_playlist() {
  if [[ -z "$4" ]]; then
    echo "Usage: $0 "rm" playlist <user_id> <playlist_id>"
    exit 1
  fi
  user_id="$3"
  playlist_id="$4"
  # check if user exists
  if [[ -z "$(./s2postgre.sh list user | grep "User: $user_id")" ]]; then
    echo "User $user_id does not exist."
    exit 1
  fi
  ./../S2_Mongo/S2_Mongo "rm" playlist "$user_id" "$playlist_id"
  exit 0
}

# remove genre from DB
function handlerm_genre() {
  if [[ -z "$3" ]]; then
    echo "Usage: $0 "rm" genre <genre_id>"
    exit 1
  fi
  genre_id="$3"
  ./s2postgre.sh "rm" genre "$genre_id"
  exit 0
}

# remove track from history
function handlerm_history() {
  if [[ -z "$5" ]]; then
    echo "Usage: $0 "rm" history <user_id> <track_id> <played_at>"
    exit 1
  fi
  user_id="$3"
  # check if user exists
  if [[ -z "$(./s2postgre.sh list user | grep "User: $user_id")" ]]; then
    echo "User $user_id does not exist."
    exit 1
  fi
  track_id="$4"
  # check if track exists
  if [[ -z "$(./s2postgre.sh get track "$track_id" | grep "ID: $track_id")" ]]; then
    echo "Track $track_id does not exist."
    exit 1
  fi
  played_at="$5"
  ./../S2_Scylla/S2_Scylla "rm" "history" "$user_id" "$track_id" "$played_at"
  exit 0
}

# remove like
function handlerm_like() {
  if [[ -z "$4" ]]; then
    echo "Usage: $0 "rm" like <user_id> <track_id>"
    exit 1
  fi
  user_id="$3"
  track_id="$4"
  # check if user exists
  if [[ -z "$(./s2postgre.sh list user | grep "User: $user_id")" ]]; then
    echo "User $user_id does not exist."
    exit 1
  fi
  # check if track exists
  if [[ -z "$(./s2postgre.sh get track "$track_id" | grep "ID: $track_id")" ]]; then
    echo "Track $track_id does not exist."
    exit 1
  fi
  ./../S2_Mongo/S2_Mongo "rm" track "$user_id" likes "$track_id"
  exit 0
}

# remove time entry
function handlerm_time() {
  if [[ -z "$4" ]]; then
    echo "Usage: $0 "rm" time <user_id> <track_id>"
    exit 1
  fi
  user_id="$3"
  # check if user exists
  if [[ -z "$(./s2postgre.sh list user | grep "User: $user_id")" ]]; then
    echo "User $user_id does not exist."
    exit 1
  fi
  track_id="$4"
  # check if track exists
  if [[ -z "$(./s2postgre.sh get track "$track_id" | grep "ID: $track_id")" ]]; then
    echo "Track $track_id does not exist."
    exit 1
  fi
  ./../S2_Scylla/S2_Scylla "rm" "time" "$user_id" "$track_id"
  exit 0
}

# list all users
function handleuser_list() {
  ./s2postgre.sh list user
  exit 0
}

# list all artists
function handleartist_list() {
  ./s2postgre.sh list artist
  exit 0
}

# list all albums
function handlealbum_list() {
  ./s2postgre.sh list album
  exit 0
}

# list tracks
function handletrack_list() {
  ./s2postgre.sh list track
  exit 0
}

# list playlists for a user
function handleplaylist_list() {
  if [[ -z "$3" ]]; then
    echo "Usage: $0 list playlist <user_id>"
    exit 1
  fi
  user_id="$3"
  # check if user exists
  if [[ -z "$(./s2postgre.sh list user | grep "User: $user_id")" ]]; then
    echo "User $user_id does not exist."
    exit 1
  fi
  ./../S2_Mongo/S2_Mongo list playlist "$user_id"
  exit 0
}

# list all genres
function handlegenre_list() {
  ./s2postgre.sh list genre
  exit 0
}

# list listening history for a user
function handlehistory_list() {
  if [[ -z "$4" ]]; then
    echo "Usage: $0 list history <user_id> <num>"
    exit 1
  fi
  user_id="$3"
  # check if user exists
  if [[ -z "$(./s2postgre.sh list user | grep "User: $user_id")" ]]; then
    echo "User $user_id does not exist."
    exit 1
  fi
  num=$4
  ./../S2_Scylla/S2_Scylla recent tracks "$user_id" $num
  exit 0
}

# list liked songs for a user
function handlelike_list() {
  if [[ -z "$4" ]]; then
    echo "Usage: $0 list like <user_id>"
    exit 1
  fi
  user_id="$3"
  # check if user exists
  if [[ -z "$(./s2postgre.sh list user | grep "User: $user_id")" ]]; then
    echo "User $user_id does not exist."
    exit 1
  fi
  ./../S2_Mongo/S2_Mongo list track "$user_id" likes
  exit 0
}



# list most played from user
function handlemostp_list() {
  if [[ -z "$4" ]]; then
    echo "Usage: $0 list mostp <user_id> <top_num>"
    exit 1
  fi
  user_id="$3"
  # check if user exists
  if [[ -z "$(./s2postgre.sh list user | grep "User: $user_id")" ]]; then
    echo "User $user_id does not exist."
    exit 1
  fi
  top_num="$4"
  ./../S2_Scylla/S2_Scylla top tracks "$user_id" "$top_num"
  exit 0
}

function handleadd() {
  for sub in "${accepted_subcommands[@]}"; do
    case "$2" in
      "user") handleadd_user "$@" ;;
      "artist") handleadd_artist "$@" ;;
      "album") handleadd_album "$@" ;;
      "track") handleadd_track "$@" ;;
      "playlist") handleadd_playlist "$@" ;;
      "genre") handleadd_genre "$@" ;;
      "history") handleadd_history "$@" ;;
      "like") handleadd_like "$@" ;;
      "time") handleadd_time "$@" ;;
      # mostp is not added, it's queried
      *) echo "Unknown subcommand for add: $2"; exit 1 ;;
    esac
  done
  exit 0
}

function handlerm() {
  for sub in "${accepted_subcommands[@]}"; do
    case "$2" in
      "user") handlerm_user "$@" ;;
      "artist") handlerm_artist "$@" ;;
      "album") handlerm_album "$@" ;;
      "track") handlerm_track "$@" ;;
      "playlist") handlerm_playlist "$@" ;;
      "genre") handlerm_genre "$@" ;;
      "history") handlerm_history "$@" ;;
      "like") handlerm_like "$@" ;;
      "time") handlerm_time "$@" ;;
      # mostp is not removed, it's queried
      *) echo "Unknown subcommand for rm: $2"; exit 1 ;;
    esac
  done
  exit 0
}

function handletime_list() {
  $userid="$3"
  $recent_num="$4"
  ../S2_Scylla/S2_Scylla recent tracks "$user_id" "$recent_num"
}

function handlelist() {
  for sub in "${accepted_subcommands[@]}"; do
    case "$2" in
      "user") handleuser_list "$@" ;;
      "artist") handleartist_list "$@" ;;
      "album") handlealbum_list "$@" ;;
      "track") handletrack_list "$@" ;;
      "playlist") handleplaylist_list "$@" ;;
      "genre") handlegenre_list "$@" ;;
      "history") handlehistory_list "$@" ;;
      "like") handlelike_list "$@" ;;
      "time") handletime_list "$@" ;;
      "mostp") handlemostp_list "$@" ;;
      *) echo "Unknown subcommand for list: $2"; exit 1 ;;
    esac
  done
  exit 0
}

function handleget() {
  # user, artist, album, track, genre, playlist
  for sub in "${accepted_subcommands[@]}"; do
    case "$2" in
      "user") ./s2postgre.sh get user "$3" ; exit 0;;
      "artist") ./s2postgre.sh get artist "$3"; exit 0 ;;
      "album") ./s2postgre.sh get album "$3"; exit 0 ;;
      "track") ./s2postgre.sh get track "$3"; exit 0 ;;
      "genre") ./s2postgre.sh get genre "$3"; exit 0 ;;
      "playlist") ./../S2_Mongo/S2_Mongo get playlist "$3" "$4"; exit 0 ;;
    esac
  done
  exit 0
}

case "$1" in
  "add") handleadd "$@" ;;
  "rm") handlerm "$@" ;;
  "list") handlelist "$@" ;;
  "get") handleget "$@" ;;
  *) echo "Unknown command: $1"; exit 1 ;;
esac
exit 0
