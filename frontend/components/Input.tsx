import React from 'react';
import {
  View,
  Text,
  TextInput,
  StyleSheet,
  TextInputProps,
  ViewStyle,
  TextStyle,
} from 'react-native';

interface InputProps extends TextInputProps {
  label?: string;
  error?: string;
  containerStyle?: ViewStyle;
  inputStyle?: TextStyle;
  leftIcon?: React.ReactNode;
}

export const Input: React.FC<InputProps> = ({
  label,
  error,
  containerStyle,
  inputStyle,
  leftIcon,
  ...rest
}) => {
  return (
    <View style={[styles.container, containerStyle]}>
      {label && <Text style={styles.label}>{label}</Text>}
      <View style={[styles.inputWrapper, error && styles.inputError]}>
        {leftIcon && <View style={styles.icon}>{leftIcon}</View>}
        <TextInput
          style={[styles.input, inputStyle]}
          placeholderTextColor="#999"
          {...rest}
        />
      </View>
      {error ? <Text style={styles.errorText}>{error}</Text> : null}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginBottom: 16,
  },
  label: {
    fontSize: 14,
    color: '#333',
    marginBottom: 6,
    fontWeight: '600',
  },
  inputWrapper: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#f8f8f8',
    paddingHorizontal: 14,
    borderRadius: 12,
    borderWidth: 1,
    borderColor: '#ddd',
  },
  input: {
    flex: 1,
    height: 48,
    fontSize: 16,
    color: '#000',
  },
  inputError: {
    borderColor: '#ff3b30',
  },
  errorText: {
    marginTop: 6,
    fontSize: 13,
    color: '#ff3b30',
  },
  icon: {
    marginRight: 10,
  },
});
