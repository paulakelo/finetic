class SmsPayload {
  const SmsPayload({
    required this.userId,
    required this.rawSms,
  });

  final String userId;
  final String rawSms;

  factory SmsPayload.fromJson(Map<String, dynamic> json) {
    return SmsPayload(
      userId: json['user_id'] as String? ?? '',
      rawSms: json['raw_sms'] as String? ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'user_id': userId,
      'raw_sms': rawSms,
    };
  }
}