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
      print('Successfully ingested M-Pesa transaction.');
    } else {
      print('Failed to sync M-Pesa transaction with backend.');
    }
  }
}