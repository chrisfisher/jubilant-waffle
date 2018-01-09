// @flow

import React, { Component } from 'react';
import { StackNavigator } from 'react-navigation';

import {
  StyleSheet,
  Text,
  View,
  FlatList,
  TouchableOpacity,
} from 'react-native';

const App = (props: {}) => {
  const NavigatorComponent = StackNavigator({
    FilmListScreen: { screen: FilmList },
    FilmScreen: { screen: Film },
  });
  return <NavigatorComponent />;  
};

export default App;

type FilmListProps = {
  navigation: any,
};

class FilmList extends Component<FilmListProps> {
  static navigationOptions = {
    title: 'Films',
  };
  
  renderItem = ({item: {title}}: any) => {
    const onPress = () => this.props.navigation.navigate('FilmScreen', {title});
    return (
      <TouchableOpacity onPress={onPress}>
        <View style={styles.item}>
          <Text style={styles.itemTitle}>{title}</Text>
        </View>
      </TouchableOpacity>
    );
  };

  render() {
    const data = [{
      key: '1',
      title: 'Film 1',
    }, {
      key: '2',
      title: 'Film 2',
    }];

    return (
      <FlatList 
        data={data}
        renderItem={this.renderItem}
        style={styles.list}
      />
    );
  }
}

class Film extends Component<{}> {
  static navigationOptions = ({navigation: {state}}) => ({
    title: state.params ? state.params.title : 'Film'
  });

  render() {
    return (
      <View style={styles.film}>
      </View>
    );
  }
}

const styles = StyleSheet.create({
  list: {
    flex: 1,
  },
  item: {
    flex: 1,
    height: 120,
    padding: 15,
  },
  itemTitle: {
    color: 'black',
  },
  film: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
});
