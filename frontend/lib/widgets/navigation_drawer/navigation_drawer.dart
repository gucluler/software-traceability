import 'package:flutter/material.dart';
import 'package:frontend/routing/route_names.dart';
import 'package:frontend/widgets/navigation_drawer/drawer_item.dart';
import 'package:frontend/widgets/navigation_drawer/navigation_drawer_header.dart';

class NavigationDrawer extends StatelessWidget {
  const NavigationDrawer({Key key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 300,
      decoration: BoxDecoration(
        color: Colors.white,
        boxShadow: [BoxShadow(color: Colors.black12, blurRadius:16),],
      ),
      child: Column(
        children: <Widget>[
          NavigationDrawerHeader(),
          // Can combine
          DrawerItem('Start', Icons.keyboard, LoginRoute),
          DrawerItem('About', Icons.help, AboutRoute),
        ],
      )
    );
  }
}