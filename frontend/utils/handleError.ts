import { Alert } from 'react-native';

export function handleError(err: unknown) {
  if (err instanceof Error) {
    // Check if it's an email confirmation error
    if (err.message.includes('Please check your email for a confirmation link')) {
      Alert.alert(
        'Email Confirmation Required',
        'Please check your email for a confirmation link. If you haven\'t received it, try signing up again.',
        [
          {
            text: 'OK',
            style: 'default',
          },
        ]
      );
    } else {
      Alert.alert('Error', err.message);
    }
  } else {
    Alert.alert('Error', 'Something went wrong');
  }
}

