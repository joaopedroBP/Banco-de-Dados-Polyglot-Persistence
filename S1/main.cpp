#include <array>
#include <cstdio>
#include <iostream>
#include <memory>
#include <sstream>
#include <string>

using namespace std;

struct FileCloser {
  void operator()(FILE *f) const noexcept {
    if (f)
      pclose(f);
  }
};

string runCommand(const string &command) {
  array<char, 128> buffer{};
  string result;

  unique_ptr<FILE, FileCloser> pipe(popen(command.c_str(), "r"));
  if (!pipe) {
    throw runtime_error("popen() failed!");
  }

  while (fgets(buffer.data(), buffer.size(), pipe.get()) != nullptr) {
    result += buffer.data();
  }
  return result;
}

/* tables
users
artists
albums
genres
songs
playlists (mongo)
history (scylla)
recommends (scylla)
wrapped (scylla) (same as spotify wrapped)
*/

int addLikesPlaylist(string username) {
  // returns the like playlist id
  string res = runCommand("s2mongo add playlist " + username + " likes");
  if (res.length() > 0)
    return atoi(res.c_str());
  return -1;
}

string createUser(string username, string email, string passwd) {
  // returns the username
  return runCommand("s2postgre.sh add user \"" + username + "\" \"" + email +
                    "\" \"" + passwd + "\"");
}

int addUser(string username, string email, string passwd) {
  // returns the like playlist id
  return addLikesPlaylist(createUser(username, email, passwd));
}

int rmUser(string username) {
  // returns 0 on success
  return atoi(runCommand("s2postgre.sh rm user " + username).c_str());
}

int addTrack(string name, string album, string duration) {
  // duration in seconds
  return atoi(runCommand("s2postgre.sh add track \"" + name + "\" \"" + album +
                         "\" " + duration)
                  .c_str());
}

int addAlbum(string name, string artist, string genre) {
  // returns album id
  return atoi(runCommand("s2postgre.sh add album \"" + name + "\" \"" + artist +
                         "\" \"" + genre + "\"")
                  .c_str());
}

int addArtist(string name) {
  // returns artist id
  return atoi(runCommand("s2postgre.sh add artist \"" + name + "\"").c_str());
}

int addPlaylist(string username, string playlistName) {
  // returns playlist id
  return atoi(runCommand("s2mongo add playlist " + username + " \"" +
                         playlistName + "\"")
                  .c_str());
}

string getTrack(string trackid) {
  // returns track id
  string res = runCommand("s2postgre.sh get track " + trackid);
  if (res.length() > 0)
    return res;
  return "";
}

int addTrackToPlaylist(string username, string playlistName, string trackId) {
  // returns 0 on success
  // sees if exists on postgres
  if (getTrack(trackId) != "")
    return atoi(runCommand("s2mongo add track \"" + username + "\" \"" +
                           playlistName + "\" " + trackId)
                    .c_str());
  return -1;
}

int rmPlaylist(string username, string playlistName) {
  // returns 0 on success
  return atoi(runCommand("s2mongo rm playlist " + username + " \"" +
                         playlistName + "\"")
                  .c_str());
}

int rmTrackFromPlaylist(string username, string playlistName, int trackId) {
  // returns 0 on success
  return atoi(runCommand("s2mongo rm track " + username + " \"" + playlistName +
                         "\" " + to_string(trackId))
                  .c_str());
}

int addToHistory(string username, int trackId) {
  // returns 0 on success
  return atoi(
      runCommand("s2scylla add history " + username + " " + to_string(trackId))
          .c_str());
}

int addToRecommends(string username, int trackId) {
  // returns 0 on success
  return atoi(runCommand("s2scylla add recommend " + username + " " +
                         to_string(trackId))
                  .c_str());
}

int addToWrapped(string username, int trackId) {
  // returns 0 on success
  return atoi(
      runCommand("s2scylla add wrapped " + username + " " + to_string(trackId))
          .c_str());
}

int addGenre(string genreName) {
  // returns genre id
  return atoi(
      runCommand("s2postgre.sh add genre \"" + genreName + "\"").c_str());
}

int addLike(string username, int trackId) {
  // returns 0 on success
  return atoi(runCommand("s2mongo add song \"" + username + "\" likes " +
                         to_string(trackId))
                  .c_str());
}

string getTrackName(string trackId) {
  // returns track name
  return runCommand("s2postgre.sh get track name " + trackId);
}

int listPlaylistSongs(string username, string playlist) {
  string trackIDs = runCommand("s2mongo list songs \"" + username + "\" \"" +
                               playlist + "\"");
  // separa trackIds e printa os nomes das musicas
  istringstream ss(trackIDs);
  string trackID;
  while (getline(ss, trackID)) {
    cout << getTrackName(trackID);
  }
  return 0;
}

int main(int argc, char *argv[]) {
  if (argc < 2) {
    return 1;
  }

  // see which command to run
  string command = argv[1];
  string table = argv[2];
  if (command == "add") {
    if (table == "user") {
      if (argc != 6) {
        return 2;
      }
      cout << addUser(argv[3], argv[4], argv[5]) << endl;
    }
  }
}
