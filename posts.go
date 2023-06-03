package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
)


type PostSearch struct {

    Id                  int64  `json:"id"`
    User                string `json:"users"`
    Icon                string `json:"icon"`
    Content             string `json:"content"`
    Like                int64  `json:"likes"`
    Dislikes            int64  `json:"dislikes"`
    Draft               bool   `json:"draft_state"`
    ReportsId           string `json:"reports_id"`
    Timestamp           string `json:"timestamp"`
    Score               int64  `json:"score"`
    Edited              bool   `json:"edited_true"`
    CommentsTrue        bool   `json:"comments_true"`
    OptionsTrue         bool   `json:"options_true"`
    LocalEdit           bool   `json:"local_edit"`
    RemovedFlag         bool   `json:"removed_flag"`
    ComsImage           ComsImage `image`
    ComsVideo           ComsVideo `video`

}

type CountsInd struct{
    Counts                  int64  `json:"count"`

}

type ComsImage struct {
    ImageTrue           bool   `json:"image_true"`
    ImageHash           string `json:"image_hash"`

}

type ComsVideo struct {
    VideoTrue           bool   `json:"video_true"`
    VideoLink           string `json:"video_link"`
}



const baseQueryPost = "SELECT  id, users, icon, content, likes, dislikes, content_type, draft, reports_id, time_stamp, score, edited_true, image_true, image_hash, video_true, video_link FROM posts "
const baseQueryCreatePost = "INSERT INTO posts (id, users, icon, content, likes, dislikes, content_type, draft, reports_id, time_stamp, score, edited_true, image_true, image_hash, video_true, video_link)"
const baseQueryAction = "SELECT id, liked_true, dislike_true FROM action_like "


func getPostsByName(user string, ctx context.Context) ([]PostSearch, error) {

    var users string
    var friends []FriendSearch
    var following []FollowingSearch
    var err error
    requester := user

    following, err = getFolByName(user, ctx)
    friends, err = getFriendByName(user, ctx)

    for i := 0; i < len(following); i++{
        users = users + (fmt.Sprintf("'%s',",following[i].Followed))
    }

    for i := 0; i < len(friends); i++{
        users = users + (fmt.Sprintf("'%s',",friends[i].Friend))
    }

	tsql := baseQueryPost + fmt.Sprintf(" WHERE content_type='POST' and users in (%s 'NEXUS','%s') ORDER BY time_stamp DESC;",users,user)

	results, err := api.Upstream.Db.QueryContext(ctx, tsql, nil)

	if err != nil {
		api.Logger.Error(fmt.Sprintf("Error executing query: %v", err))
		return nil, err
	}

	defer results.Close()

	out := make([]PostSearch, 0)

	for results.Next() {
		var f PostSearch
		var id sql.NullInt64
        var user sql.NullString
        var icon sql.NullString
        var content sql.NullString
        var like sql.NullInt64
        var dislikes sql.NullInt64
        var contentType sql.NullString
        var draft sql.NullString
        var reportsId sql.NullString
        var timestamp sql.NullString
        var score sql.NullInt64
        var editedTrue sql.NullString
        var imageTrue sql.NullString
        var imageHash sql.NullString
        var videoTrue sql.NullString
        var videoLink sql.NullString

		if err := results.Scan(
               &id,
               &user,
               &icon,
               &content,
               &like,
               &dislikes,
               &contentType,
               &draft,
               &reportsId,
               &timestamp,
               &score,
               &editedTrue,
               &imageTrue,
               &imageHash,
               &videoTrue,
               &videoLink,



        ); err != nil {
		    api.Logger.Error(fmt.Sprintf("File %s: %s", f.User, err.Error()))
		}

		if id.Valid {
			f.Id = id.Int64
		}
        if user.Valid {
           f.User = user.String
        }
        if icon.Valid {
           f.Icon = icon.String
        }
        if content.Valid {
           f.Content = content.String
        }
        if like.Valid {
           f.Like = like.Int64
        }
        if dislikes.Valid {
           f.Dislikes = dislikes.Int64
        }
        if draft.Valid {
            if  (draft.String == "true"){
                f.Draft = true
            }else{
                f.Draft = false
            }
        }

        if contentType.Valid {
           f.ContentType = contentType.String
        }
        if reportsId.Valid {
            f.ReportsId = reportsId.String
        }
        if timestamp.Valid {
             f.Timestamp = timestamp.String
        }
        if score.Valid {
             f.Score = score.Int64
        }
        str:= strconv.FormatInt(f.Id, 10)

        if imageTrue.Valid {
            if  (imageTrue.String == "true"){
                f.ComsImage.ImageTrue = true
            }else{
                f.ComsImage.ImageTrue = false
            }
        }
        if imageHash.Valid {
            f.ComsImage.ImageHash = imageHash.String
        }
        if videoTrue.Valid {
            if  (videoTrue.String == "true"){
                f.ComsVideo.VideoTrue = true
            }else{
                f.ComsVideo.VideoTrue = false
            }
        }
        if videoLink.Valid {
            f.ComsVideo.VideoLink = videoLink.String
        }

        if editedTrue.Valid {
            if editedTrue.String == "true" {
                f.Edited = true
            }else {
                f.Edited = false
            }
        }

        f.RemovedFlag = false
        f.CommentsTrue = false
        f.OptionsTrue = false
        f.LocalEdit = false
		out = append(out, f)

	}

	err = results.Err()

	if err != nil {
		api.Logger.Error(fmt.Sprintf("Error parsing results: %v", err))
		return nil, err
	}

	return out, nil
}

func getOtherPosts(user string, content string, ctx context.Context) ([]PostSearch, error) {

    var tsql string

    if content == "POST"{
        tsql = baseQueryPost + fmt.Sprintf(" WHERE users='%s' and content_type='POST' ORDER BY time_stamp DESC;",
        user,
        )

    }else if content == "ANNOUNCEMENT"{
        var users string
        following, err := getFolByName(user, ctx)

        if err != nil {
            api.Logger.Error(fmt.Sprintf("Error executing query: %v", err))
            return nil, err
        }

        for i := 0; i < len(following); i++{
                users = users + (fmt.Sprintf("'%s',",following[i].Followed))
        }

	    tsql = baseQueryPost + fmt.Sprintf(" WHERE content_type='ANNOUNCEMENT' and users in (%s 'NEXUS') ORDER BY time_stamp DESC;",
	    users,
	    )

    }else if content == "SUGGESTION"{
        tsql = baseQueryPost + (" WHERE content_type='SUGGESTION' ORDER BY likes DESC;")
        api.Logger.Info(tsql)

    }else if content == "PATCH"{
        tsql = baseQueryPost + fmt.Sprintf(" WHERE content_type='PATCH' AND users ='%s' ORDER BY time_stamp DESC;",
        user,
        )

    }else{
        api.Logger.Info("Incorrect Comparison")
        return nil, nil
    }

    requester := user
    posting := content

	results, err := api.Upstream.Db.QueryContext(ctx, tsql, nil)

	if err != nil {
		api.Logger.Error(fmt.Sprintf("Error executing query: %v", err))
		return nil, err
	}

	defer results.Close()

	out := make([]PostSearch, 0)

	for results.Next() {
        var f PostSearch
        var id sql.NullInt64
        var user sql.NullString
        var icon sql.NullString
        var content sql.NullString
        var like sql.NullInt64
        var dislikes sql.NullInt64
        var contentType sql.NullString
        var draft sql.NullString
        var reportsId sql.NullString
        var timestamp sql.NullString
        var score sql.NullInt64
        var editedTrue sql.NullString
        var imageTrue sql.NullString
        var imageHash sql.NullString
        var videoTrue sql.NullString
        var videoLink sql.NullString
        comsSearch := make([]ComsSearch, 0)



		if err := results.Scan(
               &id,
               &user,
               &icon,
               &content,
               &like,
               &dislikes,
               &contentType,
               &draft,
               &reportsId,
               &timestamp,
               &score,
               &editedTrue,
               &imageTrue,
               &imageHash,
               &videoTrue,
               &videoLink,



        ); err != nil {
		    api.Logger.Error(fmt.Sprintf("Other Post Search File %s: %s", f.User, err.Error()))
		}

		if id.Valid {
            f.Id = id.Int64
        }
        if user.Valid {
            f.User = user.String
        }
        if icon.Valid {
            f.Icon = icon.String
        }
        if content.Valid {
            f.Content = content.String
        }
        if like.Valid {
            f.Like = like.Int64
        }
        if dislikes.Valid {
            f.Dislikes = dislikes.Int64
        }
        if contentType.Valid {
            f.ContentType = contentType.String
        }

        if draft.Valid {
            if  (draft.String == "true"){
                f.Draft = true
            }else{
                f.Draft = false
            }
        }

        if reportsId.Valid {
            f.ReportsId = reportsId.String
        }
        if timestamp.Valid {
            f.Timestamp = timestamp.String
        }
        if score.Valid {
            f.Score = score.Int64
        }

        str:= strconv.FormatInt(f.Id, 10)
        comsSearch, err = getComs(str, f.ContentType, requester, ctx)
        f.ComsSearch = comsSearch
        f.ActionSearch, err = isLikedTrue(requester, str, posting, ctx)

        if imageTrue.Valid {
            if  (imageTrue.String == "true"){
                f.ComsImage.ImageTrue = true
            }else{
                f.ComsImage.ImageTrue = false
            }
        }
        if imageHash.Valid {
            f.ComsImage.ImageHash = imageHash.String
        }
        if videoTrue.Valid {
            if  (videoTrue.String == "true"){
                f.ComsVideo.VideoTrue = true
            }else{
                f.ComsVideo.VideoTrue = false
            }
        }
        if videoLink.Valid {
            f.ComsVideo.VideoLink = videoLink.String
        }

        if editedTrue.Valid {
            if editedTrue.String == "true" {
                f.Edited = true
            }else {
                f.Edited = false
            }
        }
        f.CommentsTrue = false
        f.OptionsTrue = false
        f.LocalEdit = false
        f.RemovedFlag = false

		out = append(out, f)

	}

	err = results.Err()

	if err != nil {
		api.Logger.Error(fmt.Sprintf("Error parsing results: %v", err))
		return nil, err
	}

	return out, nil
}


func getCounts(ctx context.Context) ([]CountsInd, error) {


    countSql := "SELECT COUNT(id) FROM posts;"
    countReturn, err := api.Upstream.Db.QueryContext(ctx, countSql , nil)

    if err != nil {
    	api.Logger.Error(fmt.Sprintf("Error executing query: %v", err))
    }

    defer countReturn.Close()

    out := make([]CountsInd, 0)

    for countReturn.Next() {
    		var f CountsInd
    		var counts sql.NullInt64

    		if err := countReturn.Scan(
                   &counts,

            ); err != nil {
    		    api.Logger.Error(fmt.Sprintf("File %s: %s", f.Counts, err.Error()))
    		}

    		if counts.Valid {
    			f.Counts = counts.Int64 + 1
    		}

    		out = append(out, f)
    }
    err = countReturn.Err()

    if err != nil {
    		api.Logger.Error(fmt.Sprintf("Error parsing results: %v", err))

    }

    return out, nil
}

func createPost(user PostSearch ,ctx context.Context) (PostSearch, error) {

    out, err:= getCounts(ctx)
    idString := strconv.FormatInt(out[0].Counts, 10)
    likeString := strconv.FormatInt(user.Like, 10)
    dislikeString := strconv.FormatInt(user.Dislikes, 10)
    scoreString := strconv.FormatInt(user.Score, 10)


    var comsImageTrue string
    var comsVideoTrue string
    var editedTrueString string

    if (user.ComsImage.ImageTrue == true){
        comsImageTrue = "true"
    }else if (user.ComsImage.ImageTrue == false){
        comsImageTrue = "false"
    }

    if (user.ComsVideo.VideoTrue == true){
        comsVideoTrue = "true"
    }else if (user.ComsVideo.VideoTrue == false){
        comsVideoTrue = "false"
    }

    if (user.Edited == true){
        comsVideoTrue = "true"
    }else if (user.Edited == false){
        comsVideoTrue = "false"
    }

    tsql := baseQueryCreatePost + fmt.Sprintf("VALUES( %s,'%s','%s','%s', %s , %s ,'%s','%s','%s','%s', %s,'%s','%s','%s','%s','%s');",
        idString,
        user.User,
        user.Icon,
        user.Content,
        likeString,
        dislikeString,
        user.ContentType,
        user.Draft,
        user.ReportsId,
        user.Timestamp,
        scoreString,
        editedTrueString,
        comsImageTrue,
        user.ComsImage.ImageHash,
        comsVideoTrue,
        user.ComsVideo.VideoLink,
    )

    api.Logger.Info(tsql)
    //Gets Repsponse
    _, err = api.Upstream.Db.ExecContext(ctx, tsql, nil)
    if err != nil {
    	log.Fatal(err)
    }
    result := user
        //returns a response with data from that record
    return result, err

}

func editPost(user PostSearch, ctx context.Context) (PostSearch, error) {

    idStringfy := strconv.FormatInt(user.Id, 10)
    likeStringfy := strconv.FormatInt(user.Like, 10)
    dislikeStringfy := strconv.FormatInt(user.Dislikes, 10)
    scoreStringfy := strconv.FormatInt(user.Score, 10)
    var comsImageTrue string
    var comsVideoTrue string
    var editedTrueString string

    if (user.ComsImage.ImageTrue == true){
        comsImageTrue = "true"
    }else if (user.ComsImage.ImageTrue == false){
        comsImageTrue = "false"
    }

    if (user.ComsVideo.VideoTrue == true){
        comsVideoTrue = "true"
    }else if (user.ComsVideo.VideoTrue == false){
        comsVideoTrue = "false"
    }

    if (user.Edited == true){
        comsVideoTrue = "true"
    }else if (user.Edited == false){
        comsVideoTrue = "false"
    }


    tsql := fmt.Sprintf("UPDATE posts SET users ='%s', content ='%s', likes ='%s', dislikes ='%s', content_type ='%s', draft='%s', reports_id ='%s', time_stamp ='%s', score ='%s', edited_true='%s', image_true='%s', image_hash='%s', video_true='%s', video_link='%s' WHERE id=%s; ",
        user.User,
        user.Content,
        likeStringfy,
        dislikeStringfy,
        user.ContentType,
        user.Draft,
        user.ReportsId,
        user.Timestamp,
        scoreStringfy,
        editedTrueString,
        comsImageTrue,
        user.ComsImage.ImageHash,
        comsVideoTrue,
        user.ComsVideo.VideoLink,
        idStringfy,

    )
    //Gets Repsponse
    _, err := api.Upstream.Db.ExecContext(ctx, tsql, nil)
    if err != nil {
    	log.Fatal(err)
    }
    result := user
        //returns a response with data from that record
    return result, err

}

func isLikedTrue(user string, postId string, postType string, ctx context.Context )(ActionSearch, error){
    var f ActionSearch
    tsql :=baseQueryAction+fmt.Sprintf("WHERE post_id = %s and post_type ='%s' and users ='%s';",
    postId,
    postType,
    user,
    )

    results, err := api.Upstream.Db.QueryContext(ctx, tsql, nil)

    rowCount := 0
    if err != nil {
    		api.Logger.Error(fmt.Sprintf("Error executing query: %v", err))
    		f.Id = 0
    		f.RecordExist = false
    		return f, err
    }

    for results.Next() {
        var id sql.NullInt64
    	var likedTrue sql.NullString
        var dislikeTrue sql.NullString


        if err := results.Scan(
            &id,
            &likedTrue,
            &dislikeTrue,

        ); err != nil {
        	api.Logger.Error(fmt.Sprintf("File %s: %s", f.LikedTrue, err.Error()))
        }

        if id.Valid {
            f.Id = id.Int64
        }

        if likedTrue.Valid{
        	 if  (likedTrue.String == "true"){
                f.LikedTrue = true
             }else{
                f.LikedTrue = false
             }
        }
        if dislikeTrue.Valid {
             if  (dislikeTrue.String == "true"){
                f.DislikeTrue = true
             }else{
                f.DislikeTrue = false
             }
        }
        f.RecordExist = true

        rowCount++

    }
    err = results.Err()
    if rowCount == 0 {
            f.Id = 0
            f.RecordExist = false
            f.LikedTrue = false
            f.DislikeTrue = false
            return f, nil

    }

    if err != nil {
    	api.Logger.Error(fmt.Sprintf("Error parsing results: %v", err))
    	f.Id = 0
    	return f, err
    }

    return f, nil

}

func deletePost(id string,ctx context.Context) (error) {

    // Query

    tsql :=  fmt.Sprintf("DELETE FROM posts WHERE id = %s;",
    id,)

    //Makes delete
    _, err := api.Upstream.Db.ExecContext(ctx, tsql, nil)
    if err != nil {
    	log.Fatal(err)
    }

        //returns a response with data from that record
    return err

}








