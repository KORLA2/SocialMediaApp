A simple social media backend built with Go + Gin, featuring user posts, comments, feed, follow system, and JWT-based authentication.

### How This Works

  ### User Registration & Activation
```
1. Users register with email, password, username, and level:

    1 → Normal user
    
    2 → Moderator (can edit other users’ posts)
    
    3 → Admin (can edit and delete other users’ posts)

2. After registration, users receive a UUID activation token to activate their account.

3. Only after activation can users sign in and receive a JWT token for protected routes.

```

### Authentication & Authorization
```

1. JWT-based authentication for all protected routes.

2. Role-based authorization for moderators and admins.
```

### Posts Management (CRUD)
```
 1. Users can create, read, update, and delete posts.

 2. Moderators and admins can edit or delete other users’ posts.
```

### Comments

```
1. Users can comment on posts.

2. Commenting restricted to authenticated users.
```

### User Feed

```
1. Users can see a feed of posts from their friends.

2. Feed sorted by recent posts or relevance.
   
  Search Posts

3. Search functionality based on post content.

4. Posts in the feed are ranked by post score (number of matched terms) to display the most relevant posts first.

```
### Follow / Unfollow Users

```
1. Users can follow or unfollow other users.

2. Followers’ posts appear in the user’s feed.
```

### Tech Stack
```
Backend: Go (Gin framework)

Database: PostgreSQL / MySQL (your choice)

Authentication: JWT + UUID activation tokens

API Docs: Swaggo (Swagger)

Testing: Postman
```

### API Documentation
```
Swagger UI available at: localhost:8008/swagger/index.html

Authorization: Use the Authorize button in Swagger and add your token like:

Bearer <your-jwt-token>
```
