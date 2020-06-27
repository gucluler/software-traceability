import 'package:flutter/material.dart';
import 'package:frontend/views/home/home_view.dart';
import 'package:frontend/widgets/project/project.dart';

class DevelopmentView extends StatefulWidget {
  final String projectName;
  final Project myProject;
  DevelopmentView({Key key, this.projectName, this.myProject}) : super(key: key);

  @override
  _DevelopmentViewState createState() => _DevelopmentViewState(projectName,myProject);
}

class _DevelopmentViewState extends State<DevelopmentView> {
  final String projectName;
  final Project myProject;
  _DevelopmentViewState(this.projectName, this.myProject);
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        centerTitle: true,
        title: Text(projectName+ " Development View"),
        backgroundColor: Colors.blue,
        actions: <Widget>[
          IconButton(
            icon: Icon(Icons.add), 
            onPressed: () {
              // Creates a pop up.
              showDialog(
                context: context,
                builder: (_)=> AlertDialog(
                    title: Text("Enter a name for your class"),
                    content: TextField(
                      maxLength: 30,
                    ),
                    actions: [
                      FlatButton(
                        child: Text("Confirm"),
                        onPressed: () {
                          Navigator.of(context, rootNavigator: true).pop('dialog');
                        },
                      ),
                    ],
                    elevation: 24.0,
                ),
                barrierDismissible: false,
              );
            },
            ),
          IconButton(
            icon: Icon(Icons.home), 
            onPressed: (){
              Navigator.push(
                context, 
                MaterialPageRoute(
                  builder: (context) => HomeView(),
                ),
              );
            }
          ),
          IconButton(
            icon: Icon(Icons.exit_to_app), 
            onPressed: (){
              Navigator.pop(context);
            })
        ],
      ),
      body: Center(child: Text("Hi")),
    );
  }
}