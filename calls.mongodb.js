use("quiz")

db.quizzes.findOne({ _id: ObjectId("676a008eb31cb394dcd442e8") })


use("quiz")
db.quizzes.find({}, { _id: 1 })

db.getMongo().getDBs()