# Social Network API

API for main features for a social network

# Enviroment Variables Structure

```
PORT=
MONGO_URL=
MAIN_DATABASE_NAME=
```

# Structures

- Authentication
  - Sign Up
  - Sign In
- Post
  - Upload
  - Delete
  - Like
  - Unlike
  - Comment
    - Create
    - Delete
    - Like
    - Unlike
    - Replies
      - Create
      - Delete
- Users Interaction
  - Follow
  - Unfollow
  - Block
  - Unblock
  - People Suggestions based in followings
  - Update Self Information
  
## User Structure

- Identificator
- Email
- First Name
- Last Name
- Username
- Password (Encrypted)
- Profile Picture
- Verified Account
- Followers
- Following
- Blocked Users
- Liked Posts
- Saved Posts
- Sign In Sessions Logs

## Post Structure

- Identificator
- Author Identificator
- Title
- Content
- Creation Date
- Likes
- Comments

## Comment Structure

- Identificator
- Post Identificator
- Author Identificator
- Content
- Creation Date
- Likes
- Replies

## Reply Structure

- Identificator
- Comment Indentificator
- Author Identificator
- Content
- Creation Date

## Like Structure

- Identificator of the user who liked
- Date of the like
