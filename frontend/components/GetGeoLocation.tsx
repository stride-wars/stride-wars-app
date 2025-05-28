import * as Location from 'expo-location';

const getLocation = async (): Promise<void> => {
  const { status } = await Location.requestForegroundPermissionsAsync();
  if (status !== 'granted') {
    alert('Permission denied!');
    return;
  }

  const location = await Location.getCurrentPositionAsync({
    accuracy: Location.Accuracy.High,
  });

  console.log('Latitude:', location.coords.latitude);
  console.log('Longitude:', location.coords.longitude);
};

const locationSubscription = Location.watchPositionAsync(
  { accuracy: Location.Accuracy.High },
  (location: Location.LocationObject) => {
    console.log('Updated location:', location.coords);
  }
);

// to stop watching :
// locationSubscription.remove();