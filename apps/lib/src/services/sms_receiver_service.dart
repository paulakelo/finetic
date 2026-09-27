import 'package:flutter/foundation.dart';

import '../network/api_service.dart';

class SmsReceiverService {
  final ApiService _apiService;
  final String _currentUserId;

  SmsReceiverService({
    required ApiService apiService,
    required String currentUserId,
  })  : _apiService = apiService,
        _currentUserId = currentUserId;

  Future<void> onMessageReceived(String sender, String body) async {
    final normalizedSender = sender.toUpperCase();
    if (!normalizedSender.contains('MPESA')) {
      return;
    }

    final success = await _apiService.syncSms(
      _currentUserId,
      body,
    );

    if (success) {
      debugPrint('Successfully ingested M-Pesa transaction.');
    } else {
      debugPrint('Failed to sync M-Pesa transaction with backend.');
    }
  }
}