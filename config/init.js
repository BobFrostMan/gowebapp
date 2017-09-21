//MongoDB initial database filling (just as am example)

//Usage: for MongoDB v2.6+:
//mongo localhost:27017/pet-operational path_to_this_file/init_entities.js

var before

//---------------------------------------------Functions----------------------------------------------------------------
function printAfter(collection, beforeCount){
    var after = db.getCollection(collection).count();
    var dif = after - beforeCount
    print("Documents after processing: " + after)
    if (dif > 0){
        print(dif + " documents inserted")
    }
    print("Finished '" + collection + "' collection init.")
}

function printBefore(collection){
    print("")
    print("Started '" + collection + "' collection init.")
    var before = db.getCollection(collection).count();
    print("Documents before processing: " + before)
    return before
}
//--------------------------------------------Clean up------------------------------------------------------------------
print("\nDropping existing database: " + db)
db.dropDatabase();
print("Database '" + db + "' removed successfully\n")

//--------------------------------------------Indexes-------------------------------------------------------------------
print("Creating indexes...")

print("Creating unique index for Users collection: by \"login\"")
db.Users.createIndex( { "login": 1 }, { unique: true } )

print("Creating unique index for Method collection: by \"name\"")
db.Method.createIndex( { "name": 1 }, { unique: true } )

print("Creating unique index for PermissionGroups collection: by \"name\"")
db.PermissionGroups.createIndex( { "name": 1 }, { unique: true } )

print("Creating unique index for Permissions collection: by \"name\" and \"type\"")
db.Permissions.createIndex( { "name": 1, "type" :1 }, { unique: true } )

print("Creating unique index for Project collection: by \"name\"")
db.Project.createIndex( { "name": 1 }, { unique: true } )

print("Creating unique index for Workflow collection: by \"name\"")
db.Workflow.createIndex( { "name": 1 }, { unique: true } )

print("Index creation finished successfully.")

//--------------------------------------------Workflow------------------------------------------------------------------
before = printBefore("Workflow")
db.Workflow.insert({
    "name" : "standard",
    "tasks" : {
        "bug" : {
            "new" : [ 
                "open", 
                "closed"
            ],
            "open" : [ 
                "done"
            ],
            "done" : [ 
                "reopen", 
                "closed"
            ],
            "reopen" : [ 
                "done"
            ]
        },
        "Improvement" : {
            "new" : [ 
                "open", 
                "closed"
            ],
            "open" : [ 
                "done"
            ],
            "done" : [ 
                "reopen", 
                "closed"
            ],
            "reopen" : [ 
                "done"
            ]
        }
    }
});
printAfter("Workflow", before)


//--------------------------------------------IssueFields---------------------------------------------------------------
before = printBefore("IssueFields")
db.IssueFields.insert({
    "name" : "standard",
    "tasks" : {
        "bug" : [ 
            "project", 
            "status", 
            "assignee", 
            "reporter", 
            "creationDate", 
            "history"
        ],
        "task" : [ 
            "project", 
            "status", 
            "assignee", 
            "reporter", 
            "creationDate", 
            "history"
        ],
        "sub-task" : [ 
            "project", 
            "status", 
            "assignee", 
            "reporter", 
            "creationDate", 
            "history"
        ],
        "improvement" : [ 
            "project", 
            "status", 
            "assignee", 
            "reporter", 
            "creationDate", 
            "history"
        ]
    }
});
printAfter("IssueFields", before)


//-------------------------------------------------Project--------------------------------------------------------------
before = printBefore("Project")
db.Project.insert({
    "issue_fields_name" : "standard",
    "name" : "EPM-CIT2",
    "workflow_name" : "standard"
});
printAfter("Project", before)


//------------------------------------------------Permissions-----------------------------------------------------------
before = printBefore("Permissions")
db.Permissions.insert({
                          "name" : "executeMethod",
                          "type" : "method",
                          "value" : "executeSomeMethod",
                          "read" : true,
                          "update" : false,
                          "execute" : true
                      });
db.Permissions.insert({
    "name" : "updateField",
    "type" : "field",
    "value" : "updateSomeField",
    "read" : true,
    "update" : true,
    "execute" : false
});

db.Permissions.insert({
  "name" : "readMethod",
  "type" : "method",
  "value" : "readSomeMethod",
  "read" : true,
  "update" : false,
  "execute" : false
});
printAfter("Permissions", before)


//------------------------------------------------PermissionGroups------------------------------------------------------
before = printBefore("PermissionGroups")
db.PermissionGroups.insert({
    "name" : "initial",
    "permissions" : [
        {
            "_id" : ObjectId("599ed05547384324b05de23f"),
            "name" : "readMethod",
            "type" : "method",
            "value" : "readSomeMethod",
            "read" : true,
            "update" : false,
            "execute" : false
        },
        {
            "_id" : ObjectId("599ed05647384324b05de240"),
            "name" : "executeMethod",
            "type" : "method",
            "value" : "executeSomeMethod",
            "read" : true,
            "update" : false,
            "execute" : true
        },
        {
            "_id" : ObjectId("599ed05647384324b05de241"),
            "name" : "updateField",
            "type" : "field",
            "value" : "updateSomeField",
            "read" : true,
            "update" : true,
            "execute" : false
        }
    ]
})
printAfter("PermissionGroups", before)


//-------------------------------------------------Users----------------------------------------------------------------
before = printBefore("Users")
db.Users.insert({
                    "login" : "Fluggegecheimen",
                    "name" : "The Bandit",
                    "password" : "$2a$10$nILEHcUGd/QGlSb358u5JuzZMzHJtArQXA9MoD0sxOl7jGLHHYX9y",
                    "groups" : [
                        {
                            "_id" : ObjectId("599ed05647384324b05de242"),
                            "name" : "initial",
                            "permissions" : [
                                {
                                    "_id" : ObjectId("599ed05547384324b05de23f"),
                                    "name" : "readMethod",
                                    "type" : "method",
                                    "value" : "readSomeMethod",
                                    "read" : true,
                                    "update" : false,
                                    "execute" : false
                                },
                                {
                                    "_id" : ObjectId("599ed05647384324b05de240"),
                                    "name" : "executeMethod",
                                    "type" : "method",
                                    "value" : "executeSomeMethod",
                                    "read" : true,
                                    "update" : false,
                                    "execute" : true
                                },
                                {
                                    "_id" : ObjectId("599ed05647384324b05de241"),
                                    "name" : "updateField",
                                    "type" : "field",
                                    "value" : "updateSomeField",
                                    "read" : true,
                                    "update" : true,
                                    "execute" : false
                                }
                            ]
                        }
                    ]
                });
db.Users.insert({
                    "login" : "user",
                    "name" : "Just user",
                    "password" : "$2a$10$VOzwuvQQv3JQv29ZGDBbM.yeABpWywzxVt8uCzIAwhydUnTt/7tjG",
                    "groups" : [
                        {
                            "_id" : ObjectId("599ed05647384324b05de242"),
                            "name" : "initial",
                            "permissions" : [
                                {
                                    "_id" : ObjectId("599ed05647384324b05de240"),
                                    "name" : "executeMethod",
                                    "type" : "method",
                                    "value" : "users",
                                    "read" : true,
                                    "update" : false,
                                    "execute" : true
                                },
                                {
                                    "_id" : ObjectId("599ed05647384324b05de241"),
                                    "name" : "executeMethod",
                                    "type" : "method",
                                    "value" : "createWorkflow",
                                    "read" : true,
                                    "update" : false,
                                    "execute" : true
                                },
                                {
                                    "_id" : ObjectId("599ed05647384324b05de243"),
                                    "name" : "executeMethod",
                                    "type" : "method",
                                    "value" : "createIssuefields",
                                    "read" : true,
                                    "update" : false,
                                    "execute" : true
                                },
                                {
                                    "_id" : ObjectId("599ed05647384324b05de243"),
                                    "name" : "executeMethod",
                                    "type" : "method",
                                    "value" : "createProject",
                                    "read" : true,
                                    "update" : false,
                                    "execute" : true
                                },
                                {
                                    "_id" : ObjectId("599ed05647384324b05de243"),
                                    "name" : "executeMethod",
                                    "type" : "method",
                                    "value" : "getUser",
                                    "read" : true,
                                    "update" : false,
                                    "execute" : true
                                }
                            ]
                        }
                    ]
                }
);

printAfter("Users", before)


//-------------------------------------------------Token----------------------------------------------------------------
before = printBefore("Token")
db.Token.insert(
{
    "expiration" : ISODate("2017-09-15T14:23:19.857Z"),
    "value" : "a0a5d315-f811-463c-85b1-cfd2fe6ccb38",
    "userId" : "59-no-such-user-just-example"
});
printAfter("Token", before)


//----------------------------------------------------------Method------------------------------------------------------
before = printBefore("Method")
db.Method.insert(
{
    "name" : "auth",
    "parameters" : [
        {
            "name" : "login",
            "required" : true,
            "type" : "string"
        },
        {
            "name" : "pass",
            "required" : true,
            "type" : "string"
        }
    ],
    "fsm" : {
        "states" : {
            "start" : {
                "start" : true,
                "transitions" : {
                    "start-find_user" : {
                        "to" : "find_user",
                        "guard" : {
                            "type" : "always"
                        },
                        "action" : {
                            "name" : "list",
                            "params" : {
                                "target" : "Users",
                                "fields" : [
                                    "login"
                                ],
                                "limit" : 1,
                                "save_as" : "user"
                            }
                        }
                    }
                }
            },
            "find_user" : {
                "transitions" : {
                    "find_user-user_found" : {
                        "to" : "user_found",
                        "guard" : {
                            "type" : "context",
                            "key" : "exists",
                            "value" : true
                        }
                    },
                    "find_user-user_not_found" : {
                        "to" : "user_not_found",
                        "guard" : {
                            "type" : "context",
                            "key" : "exists",
                            "value" : false
                        },
                        "action" : {
                            "name" : "set_result",
                            "params" : {
                                "response" : {
                                    "code" : 403.0,
                                    "data" : [
                                        {
                                            "name" : "message",
                                            "value" : "Provided user wasn't found"
                                        }
                                    ]
                                }
                            }
                        }
                    }
                }
            },
            "user_found" : {
                "transitions" : {
                    "user_found-pass_verified" : {
                        "to" : "pass_verified",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : false
                        },
                        "action" : {
                            "name" : "auth"
                        }
                    }
                }
            },
            "user_not_found" : {
                "transitions" : {}
            },
            "pass_verified" : {
                "transitions" : {
                    "pass_verified-find_token" : {
                        "to" : "find_token",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : false
                        },
                        "action" : {
                            "name" : "list",
                            "params" : {
                                "target" : "Token",
                                "fields" : [
                                    "userId"
                                ],
                                "where" : [
                                    {
                                        "name" : "userId",
                                        "value" : "_id",
                                        "from" : "entity"
                                    }
                                ],
                                "limit" : 1.0
                            }
                        }
                    },
                    "pass_verified-pass_failed" : {
                        "to" : "pass_failed",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : true
                        },
                        "action" : {
                            "name" : "set_result",
                            "params" : {
                                "response" : {
                                    "code" : 403.0,
                                    "data" : [
                                        {
                                            "name" : "message",
                                            "value" : "Provided password is invalid"
                                        }
                                    ]
                                }
                            }
                        }
                    }
                }
            },
            "find_token" : {
                "transitions" : {
                    "find_token-token_found" : {
                        "to" : "token_found",
                        "guard" : {
                            "type" : "context",
                            "key" : "exists",
                            "value" : true
                        },
                        "action" : {
                            "name" : "set_to_context",
                            "params" : {
                                "set" : {
                                    "override" : true,
                                    "token_end_date" : {
                                        "type" : "time",
                                        "units" : "seconds",
                                        "operation" : "add",
                                        "value" : 1800.0
                                    },
                                    "uuid" : {
                                        "type" : "uuid",
                                        "value" : "new"
                                    }
                                }
                            }
                        }
                    },
                    "find_token-token_not_found" : {
                        "to" : "token_not_found",
                        "guard" : {
                            "type" : "context",
                            "key" : "exists",
                            "value" : false
                        },
                        "action" : {
                            "name" : "set_to_context",
                            "params" : {
                                "set" : {
                                    "override" : true,
                                    "token_end_date" : {
                                        "type" : "time",
                                        "units" : "seconds",
                                        "operation" : "add",
                                        "value" : 1800.0
                                    },
                                    "uuid" : {
                                        "type" : "uuid",
                                        "value" : "new"
                                    }
                                }
                            }
                        }
                    }
                }
            },
            "token_found" : {
                "transitions" : {
                    "token_found-update_token" : {
                        "to" : "update_token",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : false
                        },
                        "action" : {
                            "name" : "update",
                            "params" : {
                                "target" : "Token",
                                "find" : {
                                    "fields" : [
                                        "userId"
                                    ],
                                    "where" : [
                                        {
                                            "name" : "userId",
                                            "value" : "_id",
                                            "from" : "last_entity"
                                        }
                                    ]
                                },
                                "update_values" : {
                                    "fields" : [
                                        "expiration",
                                        "value"
                                    ],
                                    "where" : [
                                        {
                                            "name" : "expiration",
                                            "value" : "token_end_date",
                                            "from" : "context"
                                        },
                                        {
                                            "name" : "value",
                                            "value" : "uuid",
                                            "from" : "context"
                                        }
                                    ]
                                }
                            }
                        }
                    }
                }
            },
            "token_not_found" : {
                "transitions" : {
                    "token_not_found-create_token" : {
                        "to" : "create_token",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : false
                        },
                        "action" : {
                            "name" : "create",
                            "params" : {
                                "target" : "Token",
                                "fields" : [
                                    "userId",
                                    "expiration",
                                    "value"
                                ],
                                "where" : [
                                    {
                                        "name" : "userId",
                                        "value" : "_id",
                                        "from" : "last_entity"
                                    },
                                    {
                                        "name" : "expiration",
                                        "value" : "token_end_date",
                                        "from" : "context"
                                    },
                                    {
                                        "name" : "value",
                                        "value" : "uuid",
                                        "from" : "context"
                                    }
                                ]
                            }
                        }
                    }
                }
            },
            "create_token" : {
                "transitions" : {
                    "create_token-token_created" : {
                        "to" : "token_created",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : false
                        },
                        "action" : {
                            "name" : "set_result",
                            "params" : {
                                "response" : {
                                    "code" : 200.0,
                                    "data" : [
                                        {
                                            "name" : "value",
                                            "value" : "value",
                                            "from" : "entity"
                                        },
                                        {
                                            "name" : "expiration",
                                            "value" : "expiration",
                                            "from" : "entity"
                                        },
                                        {
                                            "name" : "username",
                                            "value" : "name",
                                            "from" : "user"
                                        }
                                    ]
                                }
                            }
                        }
                    }
                }
            },
            "update_token" : {
                "transitions" : {
                    "update_token-token_updated" : {
                        "to" : "token_updated",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : false
                        },
                        "action" : {
                            "name" : "set_result",
                            "params" : {
                                "response" : {
                                    "code" : 200.0,
                                    "data" : [
                                        {
                                            "name" : "value",
                                            "value" : "uuid",
                                            "from" : "context"
                                        },
                                        {
                                            "name" : "expiration",
                                            "value" : "token_end_date",
                                            "from" : "context"
                                        },
                                        {
                                            "name" : "username",
                                            "value" : "name",
                                            "from" : "user"
                                        }
                                    ]
                                }
                            }
                        }
                    }
                }
            },
            "token_updated" : {
                "transitions" : {}
            },
            "token_created" : {
                "transitions" : {}
            },
            "pass_failed" : {
                "transitions" : {}
            }
        }
    }
});

db.Method.insert(
{
    "name" : "users",
    "parameters" : [],
    "fsm" : {
        "states" : {
            "start" : {
                "start" : true,
                "transitions" : {
                    "start-find_users" : {
                        "to" : "find_users",
                        "guard" : {
                            "type" : "always"
                        },
                        "action" : {
                            "name" : "list",
                            "params" : {
                                "target" : "Users",
                                "fields" : []
                            }
                        }
                    }
                }
            },
            "find_users" : {
                "transitions" : {
                    "find_users-result_returned" : {
                        "to" : "result_returned",
                        "guard" : {
                            "type" : "always"
                        },
                        "action" : {
                            "name" : "set_result"
                        }
                    }
                }
            },
            "result_returned" : {
                "transitions" : {}
            }
        }
    }
}
);

db.Method.insert(
{
    "name" : "createWorkflow",
    "parameters" : [
        {
            "name" : "name",
            "required" : true,
            "type" : "string"
        },
        {
            "name" : "tasks",
            "required" : true,
            "type" : "map[string]interface{}"
        },
        {
            "name" : "some_int",
            "required" : false,
            "type" : "float64"
        },
        {
            "name" : "some_bool",
            "required" : false,
            "type" : "bool"
        }
    ],
    "fsm" : {
        "states" : {
            "start" : {
                "start" : true,
                "transitions" : {
                    "start-find_workflow" : {
                        "to" : "find_workflow",
                        "guard" : {
                            "type" : "always"
                        },
                        "action" : {
                            "name" : "list",
                            "params" : {
                                "target" : "Workflow",
                                "fields" : [
                                    "name"
                                ]
                            }
                        }
                    }
                }
            },
            "find_workflow" : {
                "transitions" : {
                    "find_workflow-workflow_already_exists" : {
                        "to" : "workflow_already_exists",
                        "guard" : {
                            "type" : "context",
                            "key" : "exists",
                            "value" : true
                        },
                        "action" : {
                            "name" : "set_result",
                            "params" : {
                                "response" : {
                                    "code" : 405,
                                    "data" : [
                                        {
                                            "name" : "message",
                                            "value" : "Workflow with specified name already exists"
                                        }
                                    ]
                                }
                            }
                        }
                    },
                    "find_workflow-create_workflow" : {
                        "to" : "create_workflow",
                        "guard" : {
                            "type" : "context",
                            "key" : "exists",
                            "value" : false
                        },
                        "action" : {
                            "name" : "create",
                            "params" : {
                                "target" : "Workflow",
                                "fields" : [
                                    "name",
                                    "tasks"
                                ]
                            }
                        }
                    }
                }
            },
            "create_workflow" : {
                "transitions" : {
                    "create_workflow-workflow_created" : {
                        "to" : "result_returned",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : false
                        },
                        "action" : {
                            "name" : "set_result"
                        }
                    },
                    "create_workflow-workflow_creation_error" : {
                        "to" : "workflow_creation_error",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : true
                        },
                        "action" : {
                            "name" : "set_result"
                        }
                    }
                }
            },
            "workflow_creation_error" : {
                "transitions" : {}
            },
            "result_returned" : {
                "transitions" : {}
            },
            "workflow_already_exists" : {
                "transitions" : {}
            }
        }
    }
}
);

db.Method.insert(
{
    "name" : "createIssuefields",
    "parameters" : [
        {
            "name" : "name",
            "required" : true,
            "type" : "string"
        },
        {
            "name" : "tasks",
            "required" : true,
            "type" : "map[string]interface{}"
        }
    ],
    "fsm" : {
        "states" : {
            "start" : {
                "start" : true,
                "transitions" : {
                    "start-find_issue_fields" : {
                        "to" : "find_issue_fields",
                        "guard" : {
                            "type" : "always"
                        },
                        "action" : {
                            "name" : "list",
                            "params" : {
                                "target" : "IssueFields",
                                "fields" : [
                                    "name"
                                ]
                            }
                        }
                    }
                }
            },
            "find_issue_fields" : {
                "transitions" : {
                    "find_issue_fields-issue_fields_already_exists" : {
                        "to" : "issue_fields_already_exists",
                        "guard" : {
                            "type" : "context",
                            "key" : "exists",
                            "value" : true
                        },
                        "action" : {
                            "name" : "set_result",
                            "params" : {
                                "response" : {
                                    "code" : 405,
                                    "data" : [
                                        {
                                            "name" : "message",
                                            "value" : "IssueFields with specified name already exists"
                                        }
                                    ]
                                }
                            }
                        }
                    },
                    "find_issue_fields-create_issue_fields" : {
                        "to" : "create_issue_fields",
                        "guard" : {
                            "type" : "context",
                            "key" : "exists",
                            "value" : false
                        },
                        "action" : {
                            "name" : "create",
                            "params" : {
                                "target" : "IssueFields",
                                "fields" : [
                                    "name",
                                    "tasks"
                                ]
                            }
                        }
                    }
                }
            },
            "create_issue_fields" : {
                "transitions" : {
                    "create_issue_fields-issue_fields_created" : {
                        "to" : "issue_fields_created",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : false
                        },
                        "action" : {
                            "name" : "set_result"
                        }
                    },
                    "create_issue_fields-issue_fields_creation_error" : {
                        "to" : "issue_fields_creation_error",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : true
                        },
                        "action" : {
                            "name" : "set_result"
                        }
                    }
                }
            },
            "issue_fields_creation_error" : {
                "transitions" : {}
            },
            "issue_fields_created" : {
                "transitions" : {}
            },
            "issue_fields_already_exists" : {
                "transitions" : {}
            }
        }
    }
}
);

db.Method.insert(
{
    "name" : "createProject",
    "parameters" : [
        {
            "name" : "name",
            "required" : true,
            "type" : "string"
        },
        {
            "name" : "workflow_name",
            "required" : true,
            "type" : "string"
        },
        {
            "name" : "issue_fields_name",
            "required" : true,
            "type" : "string"
        }
    ],
    "fsm" : {
        "states" : {
            "start" : {
                "start" : true,
                "transitions" : {
                    "start-find_project" : {
                        "to" : "find_project",
                        "guard" : {
                            "type" : "always"
                        },
                        "action" : {
                            "name" : "list",
                            "params" : {
                                "target" : "Project",
                                "fields" : [
                                    "name"
                                ]
                            }
                        }
                    }
                }
            },
            "find_project" : {
                "transitions" : {
                    "find_project-project_found" : {
                        "to" : "project_found",
                        "guard" : {
                            "type" : "context",
                            "key" : "exists",
                            "value" : true
                        },
                        "action" : {
                            "name" : "set_result",
                            "params" : {
                                "response" : {
                                    "code" : 405,
                                    "data" : [
                                        {
                                            "name" : "message",
                                            "value" : "Project with specified name already exists"
                                        }
                                    ]
                                }
                            }
                        }
                    },
                    "find_project-project_not_found" : {
                        "to" : "project_not_found",
                        "guard" : {
                            "type" : "context",
                            "key" : "exists",
                            "value" : false
                        },
                        "action" : {}
                    }
                }
            },
            "project_not_found" : {
                "transitions" : {
                    "project_not_found-find_workflow" : {
                        "to" : "find_workflow",
                        "guard" : {
                            "type" : "always"
                        },
                        "action" : {
                            "name" : "list",
                            "params" : {
                                "target" : "Workflow",
                                "fields" : [
                                    "name"
                                ],
                                "where" : [
                                    {
                                        "name" : "name",
                                        "value" : "workflow_name",
                                        "from" : "context"
                                    }
                                ]
                            }
                        }
                    }
                }
            },
            "find_workflow" : {
                "transitions" : {
                    "find_workflow-workflow_not_found" : {
                        "to" : "workflow_not_found",
                        "guard" : {
                            "type" : "context",
                            "key" : "exists",
                            "value" : false
                        },
                        "action" : {
                            "name" : "set_result",
                            "params" : {
                                "response" : {
                                    "code" : 405,
                                    "data" : [
                                        {
                                            "name" : "message",
                                            "value" : "Can't create project. Workflow with specified name doesn't exist"
                                        }
                                    ]
                                }
                            }
                        }
                    },
                    "find_workflow-workflow_found" : {
                        "to" : "workflow_found",
                        "guard" : {
                            "type" : "context",
                            "key" : "exists",
                            "value" : true
                        },
                        "action" : {}
                    }
                }
            },
            "workflow_found" : {
                "transitions" : {
                    "workflow_found-find_issue_fields" : {
                        "to" : "find_issue_fields",
                        "guard" : {
                            "type" : "always"
                        },
                        "action" : {
                            "name" : "list",
                            "params" : {
                                "target" : "IssueFields",
                                "fields" : [
                                    "name"
                                ],
                                "where" : [
                                    {
                                        "name" : "name",
                                        "value" : "issue_fields_name",
                                        "from" : "context"
                                    }
                                ]
                            }
                        }
                    }
                }
            },
            "find_issue_fields" : {
                "transitions" : {
                    "find_issue_fields-issue_fields_found" : {
                        "to" : "issue_fields_found",
                        "guard" : {
                            "type" : "context",
                            "key" : "exists",
                            "value" : true
                        },
                        "action" : {}
                    },
                    "find_issue_fields-issue_fields_not_found" : {
                        "to" : "issue_fields_not_found",
                        "guard" : {
                            "type" : "context",
                            "key" : "exists",
                            "value" : false
                        },
                        "action" : {
                            "name" : "set_result",
                            "params" : {
                                "response" : {
                                    "code" : 405,
                                    "data" : [
                                        {
                                            "name" : "message",
                                            "value" : "Can't create project. IssueFields entity with specified name doesn't exist"
                                        }
                                    ]
                                }
                            }
                        }
                    }
                }
            },
            "issue_fields_found" : {
                "transitions" : {
                    "issue_fields_found-create_project" : {
                        "to" : "create_project",
                        "guard" : {
                            "type" : "always"
                        },
                        "action" : {
                            "name" : "create",
                            "params" : {
                                "target" : "Project",
                                "fields" : [
                                    "name",
                                    "workflow_name",
                                    "issue_fields_name"
                                ],
                                "where" : [
                                    {
                                        "name" : "name",
                                        "value" : "name",
                                        "from" : "context"
                                    },
                                    {
                                        "name" : "workflow_name",
                                        "value" : "workflow_name",
                                        "from" : "context"
                                    },
                                    {
                                        "name" : "issue_fields_name",
                                        "value" : "issue_fields_name",
                                        "from" : "context"
                                    }
                                ]
                            }
                        }
                    }
                }
            },
            "create_project" : {
                "transitions" : {
                    "create_project-project_created" : {
                        "to" : "project_created",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : false
                        },
                        "action" : {
                            "name" : "set_result"
                        }
                    },
                    "create_project-project_creation_error" : {
                        "to" : "project_creation_error",
                        "guard" : {
                            "type" : "context",
                            "key" : "failed",
                            "value" : true
                        },
                        "action" : {
                            "name" : "set_result"
                        }
                    }
                }
            },
            "project_creation_error" : {
                "transitions" : {}
            },
            "project_created" : {
                "transitions" : {}
            },
            "workflow_not_found" : {
                "transitions" : {}
            },
            "issue_fields_not_found" : {
                "transitions" : {}
            },
            "project_found" : {
                "transitions" : {}
            }
        }
    }
}
);

db.Method.insert(
{
    "name" : "getUser",
    "parameters" : [
        {
            "name" : "login",
            "type" : "string",
            "required" : true
        }
    ],
    "fsm" : {
        "states" : {
            "start" : {
                "start" : true,
                "transitions" : {
                    "start-find_user" : {
                        "to" : "find_user",
                        "guard" : {
                            "type" : "always"
                        },
                        "action" : {
                            "name" : "list",
                            "params" : {
                                "target" : "Users",
                                "fields" : [
                                    "login"
                                ]
                            }
                        }
                    }
                }
            },
            "find_user" : {
                "transitions" : {
                    "find_user-result_returned" : {
                        "to" : "result_returned",
                        "guard" : {
                            "type" : "always"
                        },
                        "action" : {
                            "name" : "set_result"
                        }
                    }
                }
            },
            "result_returned" : {
                "transitions" : {}
            }
        }
    }
}
);

printAfter("Method", before)




