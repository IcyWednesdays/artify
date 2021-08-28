# Artify

Artify is a program to update LIFX light colours based on your active Spotify listening session. It does this by:

- Monitoring for an active Spotify session
- Getting the current song
- Computing the 3 most prominent colours in that song's album artwork
- Setting the LIFX lights to those colours, with a few configurable patterns

## Usage

Create an app in Spotify, using `http://localhost:8080/callback` as the callback URL.

Start Artify by running:

```bash
 $ SPOTIFY_ID=<SPOTIFY_ID> SPOTIFY_SECRET=<SPOTIFY_SECRET> ARTIFY_PLAYER_DEVICE=<NAME_OF_SPOTIFY_LISTENING_DEVICE> ARTIFY_SCENE_NAME=<OPTIONAL: INSIDETOOUT|BLENDED> go run .
```

Click the link to auth your Spotify account, and that's it.

## TODO

- Grab the current state of lights when a Spotify session begins, so we can reset back to what it was when the session completes
- Async requests to bulbs
- Use subset of bulbs
- Monitor multiple users
- Enable/disable via specific scene (i.e. don't change anything unless FollowSpotify scene is enabled)
- Config file + more configurable options (bulb names, transition durations, etc)
- Tune colour conversion params?
- Dockerise?
